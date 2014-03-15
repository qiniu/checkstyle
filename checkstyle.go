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
)

type Problem struct {
	Position    *token.Position
	Description string
	SourceLine  string
	Type        ProblemType
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

func (f *file) checkFileContent() {
	if f.isTest() {
		return
	}

	if f.config.PackageName {
		//ref "http://golang.org/doc/effective_go.html#package-names"
		pkgName := f.ast.Name.Name
		var desc string
		if strings.Contains(pkgName, "_") {
			suggestName := strings.Replace(pkgName, "_", "/", -1)
			desc = "don't use an underscore in package name, " + pkgName + " should be " + suggestName
		} else if strings.ToLower(pkgName) != pkgName {
			desc = "don't use capital letters in package name: " + pkgName
		}
		if desc != "" {
			start := f.fset.Position(f.ast.Name.Pos())
			problem := Problem{Description: desc, Position: &start, Type: PackageName}
			f.problems = append(f.problems, problem)
		}
	}

	lineLimit := f.config.FunctionLine
	paramsNumLimit := f.config.ParamsNum
	resultsNumLimit := f.config.ResultsNum
	for _, v := range f.ast.Decls {
		switch v2 := v.(type) {
		case *ast.FuncDecl:
			start := f.fset.Position(v2.Pos())
			if lineLimit > 0 {
				startLine := start.Line
				endLine := f.fset.Position(v2.End()).Line
				lineCount := endLine - startLine
				if lineCount > lineLimit {
					problem := genFuncLineProblem(v2.Name.Name, lineCount, lineLimit, start)
					f.problems = append(f.problems, problem)
				}
			}

			_type := v2.Type
			params := _type.Params
			if paramsNumLimit != 0 && params.NumFields() > paramsNumLimit {
				problem := genParamsNumProblem(v2.Name.Name, params.NumFields(), paramsNumLimit, start)
				f.problems = append(f.problems, problem)
			}

			results := _type.Results
			if resultsNumLimit != 0 && results.NumFields() > resultsNumLimit {
				problem := genResultsNumProblem(v2.Name.Name, results.NumFields(), resultsNumLimit, start)
				f.problems = append(f.problems, problem)
			}
		}
	}
}
