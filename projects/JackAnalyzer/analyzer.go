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
	input := readFile(fileName)

	initialize(input)
	toks := tokenize()

	fileName = strings.Replace(fileName, ".jack", "", -1)
	writeFile(fileName+"T.xml", compileTokenList(toks))
	writeFile(fileName+".xml", compile(toks))
}

func compileTokenList(toks []token) []string {
	var tokenList []string
	tokenList = append(tokenList, "<tokens>")
	for _, tok := range toks {
		kind := string(tok.kind)
		str := "<" + kind + "> " + tok.str + " </" + kind + ">"
		tokenList = append(tokenList, str)
	}
	tokenList = append(tokenList, "</tokens>")
	return tokenList
}
