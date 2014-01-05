package checkstyle

import (
	"go/ast"
)

type Problem struct {
	Position    token.Position // position in source file
	Description string
	LineText    string // the source line
}

type Checker interface {
	Check(fileName string, src []byte) ([]Problem, error)
}

type checker struct {
}

func (c *checker) Check(fileName string, src []byte) ([]Problem, error) {
}

func New(config []byte) (Checker, error) {
	return &checker{}, nil
}
