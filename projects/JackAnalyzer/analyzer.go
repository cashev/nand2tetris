package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	name := os.Args[1]
	input := read(name)

	initialize(input)
	tokenize()
}
