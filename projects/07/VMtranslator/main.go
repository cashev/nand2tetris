package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	c_arithmetic = iota
	c_push
	c_pop
	c_label
	c_goto
	c_if
	c_function
	c_return
	c_call
)

var files []string
var commands []string
var pos = 0
var labelCounter = 0

func main() {
	fmt.Println("Hello, World!")
	name := os.Args[1]
	initialize(name)

	results := writeCommand()
	writeFile(results)
}

func initialize(name string) {
	fi, err := os.Stat(name)
	if err != nil {
		fmt.Println("error")
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		initializeDir(name)
	case mode.IsRegular():
		c := initializeFile(name)
		commands = append(commands, c...)
	}

}

func initializeDir(dirname string) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Println("error")
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.Contains(file.Name(), ".vm") {
			continue
		}
		c := initializeFile(file.Name())
		commands = append(commands, c...)
	}
}

func initializeFile(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		panic("error")
	}
	defer f.Close()

	var results []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// line = strings.Replace(line, " ", "", -1)
		slice := strings.Split(line, "//")
		if slice[0] != "" {
			files = append(files, f.Name())
			results = append(results, slice[0])
		}
	}
	return results
}

func writeFile(results []string) {
	fileName := files[0]
	fileName = strings.Replace(fileName, ".vm", ".asm", 1)
	file, err := os.Create(fileName)
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

func writeCommand() []string {
	var results []string
	for range commands {
		switch commandType() {
		case c_arithmetic:
			s := writeArithmetic()
			results = append(results, s)
		case c_push:
			fallthrough
		case c_pop:
			s := writePushPop()
			results = append(results, s)
		}
		pos++
	}
	return results
}

func commandType() int {
	c := commands[pos]
	v := strings.Split(c, " ")
	if len(v) == 1 {
		return c_arithmetic
	}
	switch command() {
	case "push":
		return c_push
	case "pop":
		return c_pop
	case "label":
		return c_label
	case "goto":
		return c_goto
	case "function":
		return c_function
	}

	return c_arithmetic
}

func command() string {
	c := commands[pos]
	return strings.Split(c, " ")[0]
}

func arg1() string {
	c := commands[pos]
	return strings.Split(c, " ")[1]
}

func arg2() string {
	c := commands[pos]
	return strings.Split(c, " ")[2]
}

func writePushPop() string {
	var ret string
	switch commandType() {
	case c_push:
		ret = writeArgs() + "\n"
		ret = ret + "D=A" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=D" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n"
	case c_pop:
		ret = "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "D=M" + "\n"
		ret = ret + writeArgs() + "\n"
		ret = ret + "M=D" + "\n"
	}
	return ret
}

func writeArgs() string {
	switch arg1() {
	case "local":
		arg := "@R"
		arg2, _ := strconv.Atoi(arg2())
		return arg + strconv.Itoa(1+arg2)
	case "argument":
		arg := "@R"
		arg2, _ := strconv.Atoi(arg2())
		return arg + strconv.Itoa(2+arg2)
	case "this":
		arg := "@R"
		arg2, _ := strconv.Atoi(arg2())
		return arg + strconv.Itoa(3+arg2)
	case "that":
		arg := "@R"
		arg2, _ := strconv.Atoi(arg2())
		return arg + strconv.Itoa(4+arg2)
	case "pointer":
		arg := "@R"
		arg2, _ := strconv.Atoi(arg2())
		return arg + strconv.Itoa(3+arg2)
	case "temp":
		arg := "@R"
		arg2, _ := strconv.Atoi(arg2())
		return arg + strconv.Itoa(5+arg2)
	case "constant":
		arg := "@" + arg2()
		return arg
	case "static":
		fileName := files[pos]
		arg := "@" + fileName + "." + arg2()
		return arg
	}
	return ""
}

func writeArithmetic() string {
	switch command() {
	case "add": // x + y
		// y
		ret := "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "D=M" + "\n"
		// x
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "A=M" + "\n"
		// x + y
		ret = ret + "M=D+M" + "\n"
		// end
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n"
		return ret
	case "sub": // x - y
		// y
		ret := "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "D=M" + "\n"
		// x
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "A=M" + "\n"
		// x - y
		ret = ret + "M=M-D" + "\n"
		// end
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n"
		return ret
	case "neg":
		// y
		ret := "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		// -y
		ret = ret + "M=-M" + "\n"
		// end
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n"
		return ret
	case "eq": // x = y
		// y
		ret := "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "D=M" + "\n"
		// x
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		//
		ret = ret + "D=M-D" + "\n"
		label := strconv.Itoa(labelCounter)
		ret = ret + "@label.true." + label + "\n"
		ret = ret + "D;JEQ" + "\n"
		// false
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=0" + "\n"
		ret = ret + "@label.end." + label + "\n"
		ret = ret + "0;JMP" + "\n"
		// true
		ret = ret + "(" + "label.true." + label + ")" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=-1" + "\n"
		ret = ret + "(" + "label.end." + label + ")" + "\n"
		// end
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n"
		labelCounter++
		return ret
	case "gt": // x > y
		// y
		ret := "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "D=M" + "\n"
		// x
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		//
		ret = ret + "D=M-D" + "\n"
		label := strconv.Itoa(labelCounter)
		ret = ret + "@label.true." + label + "\n"
		ret = ret + "D;JGT" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=0" + "\n"
		ret = ret + "@label.end." + label + "\n"
		ret = ret + "0;JMP" + "\n"
		ret = ret + "(" + "label.true." + label + ")" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=-1" + "\n"
		ret = ret + "(" + "label.end." + label + ")" + "\n"
		// end
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n"
		labelCounter++
		return ret
	case "lt": // x < y
		// y
		ret := "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "D=M" + "\n"
		// x
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		//
		ret = ret + "D=M-D" + "\n"
		label := strconv.Itoa(labelCounter)
		ret = ret + "@label.true." + label + "\n"
		ret = ret + "D;JLT" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=0" + "\n"
		ret = ret + "@label.end." + label + "\n"
		ret = ret + "0;JMP" + "\n"
		ret = ret + "(" + "label.true." + label + ")" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=-1" + "\n"
		ret = ret + "(" + "label.end." + label + ")" + "\n"
		// end
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n"
		labelCounter++
		return ret
	case "and":
		// y
		ret := "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "D=M" + "\n"
		// x
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		// and
		ret = ret + "M=D&M" + "\n"
		// end
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n"
		return ret
	case "or":
		// y
		ret := "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "D=M" + "\n"
		// x
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		// or
		ret = ret + "M=D|M" + "\n"
		// end
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n"
		return ret
	case "not":
		// y
		ret := "@R0" + "\n"
		ret = ret + "M=M-1" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=!M" + "\n"
		// end
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n"
		return ret
	}

	return ""
}
