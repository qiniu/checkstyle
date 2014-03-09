package testdata

import (
	"fmt"
)

func hello(a, b int) (int, int, int, int) {
	fmt.Println("hello")
}

func hello2(a, b int, c ...int) (int, int, int) {
	fmt.Println("... args")
}

type z int

func (_ *z) hello3(a, b, c int) (int, int) {
	fmt.Println("receiver")
}
