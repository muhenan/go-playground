package main

import (
	"fmt"
	"go-playground/calculator"
	"go-playground/packages/concurrent"
	"go-playground/packages/slice"
	"go-playground/webapp"
)

func main() {
	fmt.Println("Hello, World!")
	var n = 3
	switch n {
	case 0:
		calculator.Calculator()
	case 1:
		slice.Slice()
	case 2:
		concurrent.Concurrent()
	case 3:
		webapp.StartServer()
	}
}
