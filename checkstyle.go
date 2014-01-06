package checkstyle

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type Problem struct {
	Position    *token.Position
	Description string
	SourceLine  string
}

type Checker interface {
	Check(fileName string, src []byte) ([]Problem, error)
}

type checker struct {
	FunctionComment bool
	FileLine        int
	FunctionLine    int
	MaxIndent       int
	IndentFormat    bool
}

func New(config []byte) (Checker, error) {
	return &checker{FileLine: 10}, nil
}

func (c *checker) Check(fileName string, src []byte) (ps []Problem, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return (&file{fileName, src, c, f, fset, []Problem{}}).check(), nil
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
	f.checkFileLine()
	return f.problems
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
		fmt.Println("file lines", lineCount)
		if lineCount > lineLimit {
			desc := strconv.Itoa(lineCount) + " lines more than " + strconv.Itoa(lineLimit)
			problem := Problem{Description: desc}
			f.problems = append(f.problems, problem)
			return false
		}
		return true
	})
}
