package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")

	dirName := os.Args[1]
	files := readDir(dirName)

	for _, f := range files {
		analyzeFile(f)
	}

}

func analyzeFile(fileName string) {
	fmt.Println(fileName)
	input := readFile(fileName)

	initialize(input)
	toks := tokenize()
	// compile(toks)

	fileName = strings.Replace(fileName, ".jack", "", -1)
	tResults := compileTokenList(toks)
	writeFile(fileName+"T.xml", tResults)
}
