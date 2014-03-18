package testdata

import "fmt"

const (
	Abc  = 0
	_Abc = 0
)

const _AB = 0

var _Bc = 0

var (
	_CdE = 0
	_DE  = 0
)

type _Ac struct {
	AAAA   int
	_BBBB  int
	_AAAAa int
}

type x interface {
}

func (xz *_Ac) HelloWorld(ab, _ab int) (bc int) {
	var _xyz int
	fmt.Println(_xyz)
	return
}

func (_ *_Ac) helloWorld(ab, _ab int) (bc int) {
	var _xyz int
	fmt.Println(_xyz)
	return
}
