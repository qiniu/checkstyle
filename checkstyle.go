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
	Formated     ProblemType = "formated"
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
	f.checkFunctionLine()
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
	if bytes.Compare(src, f.src) != 0 {
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
			return false
		}
		return true
	})
}

func (f *file) checkFunctionLine() {
	if f.isTest() {
		return
	}

	lineLimit := f.config.FunctionLine

	if lineLimit == 0 {
		return
	}
	for _, v := range f.ast.Decls {
		switch v := v.(type) {
		case *ast.FuncDecl:
			start := f.fset.Position(v.Pos())
			startLine := start.Line
			endLine := f.fset.Position(v.End()).Line
			lineCount := endLine - startLine
			if lineCount > lineLimit {
				desc := "func " + v.Name.Name + "() " + strconv.Itoa(lineCount) + " lines more than " + strconv.Itoa(lineLimit)
				problem := Problem{Description: desc, Position: &start, Type: FunctionLine}
				f.problems = append(f.problems, problem)
			}
		}
	}
}
