package testdata

import "fmt"

func HelloWorld() {
	defer fmt.Println("world")

	fmt.Println("hello")
}
