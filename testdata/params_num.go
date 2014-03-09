package testdata

import (
	"fmt"
)

func hello(a, b, c, _ int) {
	fmt.Println("hello")
}

func hello2(a, b int, c ...int) {
	fmt.Println("... args")
}

type z int

func (_ *z) hello(a, b int) {
	fmt.Println("one receiver")
}
