package checkstyle

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
	"testing"
)

const baseDir = "testdata/"

func readFile(fileName string) []byte {
	file, _ := ioutil.ReadFile(baseDir + fileName)
	return file
}

func TestFileLine(t *testing.T) {
	fileName := "fileline.go"
	file := readFile(fileName)
	_checkerOk := checker{FileLine: 9}
	ps, err := _checkerOk.Check(fileName, file)
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 0 {
		t.Fatal("expect no error")
	}

	_checkerFail := checker{FileLine: 8}
	ps, _ = _checkerFail.Check(fileName, file)
	if len(ps) != 1 || ps[0].Type != FileLine {
		t.Fatal("expect an error")
	}

	//pos is at file end
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, fileName, file, parser.ParseComments)
	if reflect.DeepEqual(ps[0], fset.Position(f.End())) {
		t.Fatal("file line problem position not match")
	}
}

func TestFunctionLine(t *testing.T) {
	fileName := "functionline.go"
	file := readFile(fileName)
	_checkerOk := checker{FunctionLine: 9}
	ps, err := _checkerOk.Check(fileName, file)
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 0 {
		t.Fatal("expect no error")
	}

	_checkerFail := checker{FunctionLine: 8}
	ps, _ = _checkerFail.Check(fileName, file)
	if len(ps) != 1 || ps[0].Type != FunctionLine {
		t.Fatal("expect an error")
	}

	if ps[0].Position.Filename != fileName {
		t.Fatal("file name is not correct")
	}

	if ps[0].Position.Line != 7 {
		t.Fatal("start position is not correct")
	}
}

func TestFormated(t *testing.T) {
	fileName := "formated.go"
	file := readFile(fileName)
	_checker := checker{Formated: true}
	ps, err := _checker.Check(fileName, file)
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 0 {
		t.Fatal("expect no error")
	}

	fileName = "unformated.go"
	file = readFile(fileName)
	ps, _ = _checker.Check(fileName, file)
	if len(ps) != 1 || ps[0].Type != Formated {
		t.Fatal("expect an error")
	}

	//pos is at file begin
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, fileName, file, parser.ParseComments)
	if reflect.DeepEqual(ps[0], fset.Position(f.Pos())) {
		t.Fatal("file line problem position not match")
	}
}
