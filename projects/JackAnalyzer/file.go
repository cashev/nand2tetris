package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func read(name string) []string {
	fi, err := os.Stat(name)
	if err != nil {
		panic("read error")
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return readDir(name)
	case mode.IsRegular():
		return readFile(name)
	}
	panic("read Error")
}

func readDir(dirname string) []string {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic("readDir Error")
	}
	var results []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.Contains(file.Name(), ".jack") {
			continue
		}
		c := readFile(dirname + "/" + file.Name())
		results = append(results, c...)
	}
	return results
}

func readFile(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		panic("error")
	}
	defer f.Close()

	var results []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		results = append(results, line)
	}
	results = removeComment(results)
	return results
}

func removeComment(input []string) []string {
	// /** が始まったか判定する
	startComment := false
	hasCommentStart := func(line string) bool {
		return strings.Contains(line, COMMENTSTART) || strings.Contains(line, APICOMMENTSTART)
	}
	hasCommentEnd := func(line string) bool {
		return strings.Contains(line, COMMENTEND)
	}

	var results []string
	for _, line := range input {
		if hasCommentStart(line) {
			startComment = true
			str := strings.Split(line, COMMENTSTART)[0]
			if str != "" {
				results = append(results, str)
			}
		}
		if hasCommentEnd(line) {
			startComment = false
			str := strings.Split(line, COMMENTEND)[1]
			if str != "" {
				results = append(results, str)
			}
			continue
		}
		if startComment {
			continue
		}
		str := strings.Split(line, COMMENT)[0]
		if str != "" {
			results = append(results, str)
		}
	}
	return results
}

func writeFile(filename string, results []string) {
	outfilename := ""
	fi, _ := os.Stat(filename)
	switch mode := fi.Mode(); {
	case mode.IsDir():
		outfilename = filename + "/" + fi.Name() + ".asm"
	case mode.IsRegular():
		outfilename = strings.Replace(filename, ".vm", ".asm", 1)
	}
	file, err := os.Create(outfilename)
	if err != nil {
		fmt.Println("error")
	}
	defer file.Close()

	src := ""
	for _, s := range results {
		src = src + s
	}
	file.WriteString(src)
}
