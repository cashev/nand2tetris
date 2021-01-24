package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		slice := strings.Split(line, "//")
		if slice[0] != "" {
			_, file := filepath.Split(filename)
			file = strings.Split(file, ".")[0]
			results = append(results, slice[0])
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
