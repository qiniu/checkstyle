package checkstyle

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type ProblemType string

const (
	FileLine     ProblemType = "file_line"
	FunctionLine ProblemType = "func_line"
	ParamsNum    ProblemType = "params_num"
	ResultsNum   ProblemType = "results_num"
	Formated     ProblemType = "formated"
	PackageName  ProblemType = "pkg_name"
	CamelName    ProblemType = "camel_name"
)

type Problem struct {
	Position    *token.Position
	Description string
	// SourceLine  string
	Type ProblemType
}

type Checker interface {
	Check(fileName string, src []byte) ([]Problem, error)
	IsFatal(p *Problem) bool
}

type checker struct {
	FunctionComment bool     `json:"func_comment"`
	FileLine        int      `json:"file_line"`
	FunctionLine    int      `json:"func_line"`
	MaxIndent       int      `json:"max_indent"`
	Formated        bool     `json:"formated"`
	Fatal           []string `json:"fatal"`
	ParamsNum       int      `json:"params_num"`
	ResultsNum      int      `json:"results_num"`
	PackageName     bool     `json:"pkg_name"`
	CamelName       bool     `json:"camel_name"`
}

func New(config []byte) (Checker, error) {
	var _checker checker
	err := json.Unmarshal(config, &_checker)
	if err != nil {
		return nil, err
	}
	return &_checker, nil
}

func (c *checker) Check(fileName string, src []byte) (ps []Problem, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return (&file{fileName, src, c, f, fset, []Problem{}}).check(), nil
}

func (c *checker) IsFatal(p *Problem) bool {
	for _, v := range c.Fatal {
		if v == string(p.Type) {
			return true
		}
	}
	return false
}

type file struct {
	fileName string
	src      []byte

	config *checker

	ast  *ast.File
	fset *token.FileSet

	problems []Problem
}

func (f *file) isTest() bool {
	return strings.HasSuffix(f.fileName, "_test.go")
}

func (f *file) check() (ps []Problem) {
	f.checkFormat()
	f.checkFileLine()
	f.checkFileContent()
	return f.problems
}

func (f *file) checkFormat() {
	if !f.config.Formated {
		return
	}
	src, err := format.Source(f.src)
	if err != nil {
		panic(f.fileName + err.Error())
	}
	if len(src) != len(f.src) || bytes.Compare(src, f.src) != 0 {
		desc := "source is not formated"
		pos := f.fset.Position(f.ast.Pos())
		problem := Problem{Description: desc, Position: &pos, Type: Formated}
		f.problems = append(f.problems, problem)
	}
}

func (f *file) checkFileLine() {
	if f.isTest() {
		return
	}

	lineLimit := f.config.FileLine
	if lineLimit == 0 {
		return
	}

	f.fset.Iterate(func(_file *token.File) bool {
		lineCount := _file.LineCount()
		if lineCount > lineLimit {
			desc := strconv.Itoa(lineCount) + " lines more than " + strconv.Itoa(lineLimit)
			pos := f.fset.Position(f.ast.End())
			problem := Problem{Description: desc, Position: &pos, Type: FileLine}
			f.problems = append(f.problems, problem)
		}
		return true
	})
}

func genFuncLineProblem(name string, lineCount, lineLimit int, start token.Position) Problem {
	desc := "func " + name + "() body lines num " + strconv.Itoa(lineCount) +
		" more than " + strconv.Itoa(lineLimit)
	return Problem{Description: desc, Position: &start, Type: FunctionLine}
}

func genParamsNumProblem(name string, paramsNum, limit int, start token.Position) Problem {
	desc := "func " + name + "() params num " + strconv.Itoa(paramsNum) +
		"  more than " + strconv.Itoa(limit)
	return Problem{Description: desc, Position: &start, Type: ParamsNum}
}

func genResultsNumProblem(name string, resultsNum, limit int, start token.Position) Problem {
	desc := "func " + name + "() results num " + strconv.Itoa(resultsNum) +
		"  more than " + strconv.Itoa(limit)
	return Problem{Description: desc, Position: &start, Type: ResultsNum}
}

func (f *file) checkPkgName(pkg *ast.Ident) {
	//ref "http://golang.org/doc/effective_go.html#package-names"
	pkgName := pkg.Name
	var desc string
	if strings.Contains(pkgName, "_") {
		suggestName := strings.Replace(pkgName, "_", "/", -1)
		desc = "don't use an underscore in package name, " + pkgName + " should be " + suggestName
	} else if strings.ToLower(pkgName) != pkgName {
		desc = "don't use capital letters in package name: " + pkgName
	}
	if desc != "" {
		start := f.fset.Position(pkg.Pos())
		problem := Problem{Description: desc, Position: &start, Type: PackageName}
		f.problems = append(f.problems, problem)
	}
}

func (f *file) checkFunctionParams(fType *ast.FuncType, funcName string) {
	paramsNumLimit := f.config.ParamsNum
	resultsNumLimit := f.config.ResultsNum
	params := fType.Params
	if params != nil {
		if paramsNumLimit != 0 && params.NumFields() > paramsNumLimit {
			start := f.fset.Position(params.Pos())
			problem := genParamsNumProblem(funcName, params.NumFields(), paramsNumLimit, start)
			f.problems = append(f.problems, problem)
		}
		for _, v := range params.List {
			for _, pName := range v.Names {
				f.checkName(pName, "param", true)
			}
		}
	}

	results := fType.Results
	if results != nil {
		if resultsNumLimit != 0 && results != nil && results.NumFields() > resultsNumLimit {
			start := f.fset.Position(results.Pos())
			problem := genResultsNumProblem(funcName, results.NumFields(), resultsNumLimit, start)
			f.problems = append(f.problems, problem)
		}

		for _, v := range results.List {
			for _, rName := range v.Names {
				f.checkName(rName, "return param", true)
			}
		}
	}
}

func (f *file) checkFunctionLine(funcDecl *ast.FuncDecl) {
	lineLimit := f.config.FunctionLine
	if lineLimit <= 0 {
		return
	}
	start := f.fset.Position(funcDecl.Pos())

	startLine := start.Line
	endLine := f.fset.Position(funcDecl.End()).Line
	lineCount := endLine - startLine
	if lineCount > lineLimit {
		problem := genFuncLineProblem(funcDecl.Name.Name, lineCount, lineLimit, start)
		f.problems = append(f.problems, problem)
	}
}

func (f *file) checkFunctionDeclare(funcDecl *ast.FuncDecl) {
	f.checkFunctionLine(funcDecl)
	f.checkName(funcDecl.Name, "func", false)
	f.checkFunctionParams(funcDecl.Type, funcDecl.Name.Name)
	receiver := funcDecl.Recv
	if receiver != nil {
		f.checkName(receiver.List[0].Names[0], "receiver", true)
	}
}

func trimUnderscorePrefix(name string) string {
	if name[0] == '_' {
		return name[1:]
	}
	return name
}

func (f *file) checkName(id *ast.Ident, kind string, notFirstCap bool) {
	if !f.config.CamelName {
		return
	}
	name := trimUnderscorePrefix(id.Name)
	if name == "" {
		return
	}
	start := f.fset.Position(id.Pos())

	if strings.Contains(name, "_") {
		desc := "don't use non-prefix underscores in " + kind + " name: " + id.Name + ", please use camel name"
		problem := Problem{Description: desc, Position: &start, Type: CamelName}
		f.problems = append(f.problems, problem)
	} else if len(name) >= 5 && strings.ToUpper(name) == name {
		desc := "don't use all captial letters in " + kind + " name: " + id.Name + ", please use camel name"
		problem := Problem{Description: desc, Position: &start, Type: CamelName}
		f.problems = append(f.problems, problem)
	} else if notFirstCap && name[0:1] == strings.ToUpper(name[0:1]) {
		desc := "in function ,don't use first captial letter in " + kind + " name: " + id.Name + ", please use small letter"
		problem := Problem{Description: desc, Position: &start, Type: CamelName}
		f.problems = append(f.problems, problem)
	}
}

func (f *file) checkStruct(st *ast.StructType) {
	if st.Fields == nil {
		return
	}
	for _, v := range st.Fields.List {
		for _, v2 := range v.Names {
			f.checkName(v2, "struct field", false)
		}
	}
}

func (f *file) checkInterface(it *ast.InterfaceType) {
	if it.Methods == nil {
		return
	}
	for _, v := range it.Methods.List {
		for _, v2 := range v.Names {
			f.checkName(v2, "interface method", false)
		}
		if v3, ok := v.Type.(*ast.FuncType); ok {
			f.checkFunctionParams(v3, v.Names[0].Name)
		}
	}
}

func (f *file) checkValueName(decl *ast.GenDecl, kind string, top bool) {
	for _, spec := range decl.Specs {
		if vSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range vSpec.Names {
				f.checkName(name, kind, !top)
			}
		} else if tSpec, ok := spec.(*ast.TypeSpec); ok {
			f.checkName(tSpec.Name, kind, false)
			ast.Inspect(tSpec.Type, func(node ast.Node) bool {
				switch decl2 := node.(type) {
				case *ast.GenDecl:
					f.checkGenDecl(decl2, false)
				case *ast.FuncDecl:
					f.checkFunctionDeclare(decl2)
				case *ast.StructType:
					f.checkStruct(decl2)
				case *ast.InterfaceType:
					f.checkInterface(decl2)
				}
				return true
			})
		} else if iSpec, ok := spec.(*ast.ImportSpec); ok && iSpec.Name != nil {
			f.checkName(iSpec.Name, "import", true)
		}
	}
}

func (f *file) checkGenDecl(decl *ast.GenDecl, top bool) {
	if decl.Tok == token.CONST {
		f.checkValueName(decl, "const", top)
	} else if decl.Tok == token.VAR {
		f.checkValueName(decl, "var", top)
	} else if decl.Tok == token.TYPE {
		f.checkValueName(decl, "type", top)
	} else if decl.Tok == token.IMPORT {
		f.checkValueName(decl, "import", true)
	}
}

func (f *file) checkAssign(assign *ast.AssignStmt) {
	if assign.Tok != token.DEFINE {
		return
	}

	for _, v2 := range assign.Lhs {
		if assignName, ok := v2.(*ast.Ident); ok {
			f.checkName(assignName, "var", true)
		}
	}
}

func (f *file) checkFileContent() {
	if f.isTest() {
		return
	}

	if f.config.PackageName {
		f.checkPkgName(f.ast.Name)
	}

	for _, v := range f.ast.Decls {
		switch decl := v.(type) {
		case *ast.FuncDecl:
			f.checkFunctionDeclare(decl)
			ast.Inspect(decl.Body, func(node ast.Node) bool {
				switch decl2 := node.(type) {
				case *ast.GenDecl:
					f.checkGenDecl(decl2, false)
				case *ast.FuncDecl:
					f.checkFunctionDeclare(decl2)
				case *ast.AssignStmt:
					f.checkAssign(decl2)
				case *ast.StructType:
					f.checkStruct(decl2)
				}
				return true
			})
		case *ast.GenDecl:
			f.checkGenDecl(decl, true)
		}
	}
}
