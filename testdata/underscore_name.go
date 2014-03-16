package testdata

import (
	Fmt "fmt"
)

const (
	A_B_C = 0
)

const A_B = 0

var B_C = 0

var (
	_C_D = 0
	_D_E = 0
	__   = 0
	__a  = 0
)

type A_C struct {
	AAAAA int
	_BBBB int
}

type I_A interface {
	F_A(P_A int) (R_A int)
	F_B(int) int
}

func (Z *A_C) Hello_World(_B, A int) (_B_C, B int) {
	var X_Y_Z, _C int
	_C = 1
	Fmt.Println(X_Y_Z, __C)
	Zx := struct {
		S_A int
	}{A_B}
	Fmt.Println(Zx)
	F := func() {
		TX, TY := 0
		Fmt.Println("TX")
	}
	for I := 0; i < 1; i++ {
	}
	if X := _C; X == 1 {
	}
	F()
	return
}
