package checkstyle

import (
	"go/token"
)

type Problem struct {
	Position    token.Position
	Description string
	SourceLine  string
}

type Checker interface {
	Check(fileName string, src []byte) ([]Problem, error)
}

type checker struct {
}

func (c *checker) Check(fileName string, src []byte) ([]Problem, error) {
	return nil, nil
}

func New(config []byte) (Checker, error) {
	return &checker{}, nil
}
