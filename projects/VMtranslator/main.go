package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	writeFile(name, results)
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
			_, file := filepath.Split(filename)
			file = strings.Split(file, ".")[0]
			files = append(files, file)
			results = append(results, slice[0])
		}
	}
	return results
}

func writeFile(filename string, results []string) {
	outfilename := strings.Replace(filename, ".vm", ".asm", 1)
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

func writeCommand() []string {
	var results []string
	for range commands {
		results = append(results, "// "+commands[pos]+"\n")
		switch commandType() {
		case c_arithmetic:
			s := writeArithmetic()
			results = append(results, s)
		case c_push:
			fallthrough
		case c_pop:
			s := writePushPop()
			results = append(results, s)
		case c_label:
			s := writeLabel()
			results = append(results, s)
		case c_if:
			s := writeIf()
			results = append(results, s)
		case c_goto:
			s := writeGoto()
			results = append(results, s)
		case c_call:
			s := writeCall()
			results = append(results, s)
		case c_return:
			s := writeReturn()
			results = append(results, s)
		case c_function:
			s := writeFunction()
			results = append(results, s)
		}
		pos++
	}
	return results
}

func commandType() int {
	switch command() {
	case "push":
		return c_push
	case "pop":
		return c_pop
	case "label":
		return c_label
	case "goto":
		return c_goto
	case "if-goto":
		return c_if
	case "call":
		return c_call
	case "return":
		return c_return
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
		if arg1() == "constant" {
			ret = ret + "@" + arg2() + "\n"
			ret = ret + "D=A" + "\n"
			ret = ret + "@R0" + "\n"
			ret = ret + "A=M" + "\n"
			ret = ret + "M=D" + "\n"
			ret = ret + "@R0" + "\n"
			ret = ret + "M=M+1" + "\n" // SPの更新
			return ret
		}
		if arg1() == "static" {
			fileName := files[pos]
			arg := "@" + fileName + "." + arg2()
			ret = ret + arg + "\n"
			ret = ret + "D=M" + "\n"
			ret = ret + "@R0" + "\n"
			ret = ret + "A=M" + "\n"
			ret = ret + "M=D" + "\n"
			ret = ret + "@R0" + "\n"
			ret = ret + "M=M+1" + "\n" // SPの更新
			return ret
		}
		if arg1() == "pointer" {
			// 格納元
			arg, _ := strconv.Atoi(arg2())
			arg = 3 + arg
			ret = ret + "@" + strconv.Itoa(arg) + "\n"
			ret = ret + "D=M" + "\n"
			// 格納
			ret = ret + "@R0" + "\n"
			ret = ret + "A=M" + "\n"
			ret = ret + "M=D" + "\n"
			ret = ret + "@R0" + "\n"
			ret = ret + "M=M+1" + "\n" // SPの更新
			return ret
		}
		if arg1() == "temp" {
			// 格納元
			arg, _ := strconv.Atoi(arg2())
			arg = 5 + arg
			ret = ret + "@" + strconv.Itoa(arg) + "\n"
			ret = ret + "D=M" + "\n"
			// 格納
			ret = ret + "@R0" + "\n"
			ret = ret + "A=M" + "\n"
			ret = ret + "M=D" + "\n"
			ret = ret + "@R0" + "\n"
			ret = ret + "M=M+1" + "\n" // SPの更新
			return ret
		}

		// 格納元のアドレス解決
		ret = ret + popArg2() + "\n"
		ret = ret + "D=A" + "\n"
		ret = ret + popArg1() + "\n"
		ret = ret + "M=D+M" + "\n"
		// pushする値を保持
		ret = ret + popArg1() + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "D=M" + "\n"
		// push
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=D" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n" // SPの更新
		// 格納元のアドレス解決
		ret = ret + popArg2() + "\n"
		ret = ret + "D=A" + "\n"
		ret = ret + popArg1() + "\n"
		ret = ret + "M=M-D" + "\n"
	case c_pop:
		if arg1() == "static" {
			fileName := files[pos]
			arg := "@" + fileName + "." + arg2()
			ret = ret + "@R0" + "\n"
			ret = ret + "M=M-1" + "\n"
			ret = ret + "A=M" + "\n"
			ret = ret + "D=M" + "\n"
			ret = ret + arg + "\n"
			ret = ret + "M=D" + "\n"
			return ret
		}
		if arg1() == "pointer" {
			ret = ret + "@R0" + "\n"
			ret = ret + "M=M-1" + "\n" // SPの更新
			ret = ret + "A=M" + "\n"
			ret = ret + "D=M" + "\n" // popした値を保持
			// 格納先
			arg, _ := strconv.Atoi(arg2())
			arg = 3 + arg
			ret = ret + "@" + strconv.Itoa(arg) + "\n"
			ret = ret + "M=D" + "\n"
			return ret
		}
		if arg1() == "temp" {
			ret = ret + "M=M-1" + "\n" // SPの更新
			ret = ret + "A=M" + "\n"
			ret = ret + "D=M" + "\n" // popした値を保持
			// 格納先
			arg, _ := strconv.Atoi(arg2())
			arg = 5 + arg
			ret = ret + "@" + strconv.Itoa(arg) + "\n"
			ret = ret + "M=D" + "\n"
			return ret
		}

		// 格納先のアドレス解決
		ret = ret + popArg2() + "\n"
		ret = ret + "D=A" + "\n"
		ret = ret + popArg1() + "\n"
		ret = ret + "M=D+M" + "\n"
		// popした値
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M-1" + "\n" // SPの更新
		ret = ret + "A=M" + "\n"
		ret = ret + "D=M" + "\n" // popした値を保持
		// 値を格納
		ret = ret + popArg1() + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=D" + "\n"
		// 格納先のアドレス解決
		ret = ret + popArg2() + "\n"
		ret = ret + "D=A" + "\n"
		ret = ret + popArg1() + "\n"
		ret = ret + "M=M-D" + "\n"
	}
	return ret
}

func popArg1() string {
	var ret string
	switch arg1() {
	case "local":
		ret = "@R1"
	case "argument":
		ret = "@R2"
	case "this":
		ret = "@R3"
	case "that":
		ret = "@R4"
	case "static":
		fileName := files[pos]
		ret = "@" + fileName + "." + arg2()
	}
	return ret
}

func popArg2() string {
	return "@" + arg2()
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

func writeInit() string {
	return ""
}
func writeLabel() string {
	ret := ""
	ret = ret + "(" + arg1() + ")" + "\n"
	return ret
}
func writeGoto() string {
	ret := ""
	ret = ret + "@" + arg1() + "\n"
	ret = ret + "0;JMP" + "\n"
	return ret
}
func writeIf() string {
	ret := ""
	ret = ret + "@R0" + "\n"
	ret = ret + "M=M-1" + "\n"
	ret = ret + "A=M" + "\n"
	ret = ret + "D=M" + "\n"
	ret = ret + "@" + arg1() + "\n"
	ret = ret + "D;JNE" + "\n"
	return ret
}
func writeCall() string {
	return ""
}
func writeReturn() string {
	ret := ""
	// FRAME = LCL
	ret = ret + "@R1" + "\n"
	ret = ret + "D=M" + "\n"
	ret = ret + "@R5" + "\n"
	ret = ret + "M=D" + "\n"
	// *ARG = pop()
	ret = ret + "@R0" + "\n"
	ret = ret + "M=M-1" + "\n"
	ret = ret + "@R0" + "\n"
	ret = ret + "A=M" + "\n"
	ret = ret + "D=M" + "\n"
	ret = ret + "@R2" + "\n"
	ret = ret + "A=M" + "\n"
	ret = ret + "M=D" + "\n"
	// SP = ARG+1
	ret = ret + "@R2" + "\n"
	ret = ret + "D=M" + "\n"
	ret = ret + "@R0" + "\n"
	ret = ret + "M=D+1" + "\n"
	// THAT = *(FRAME-1)
	ret = ret + "@1" + "\n"
	ret = ret + "D=A" + "\n"
	ret = ret + "@R5" + "\n"
	ret = ret + "A=M-D" + "\n"
	ret = ret + "D=M" + "\n"
	ret = ret + "@R4" + "\n"
	ret = ret + "M=D" + "\n"
	// THIS = *(FRAME-2)
	ret = ret + "@2" + "\n"
	ret = ret + "D=A" + "\n"
	ret = ret + "@R5" + "\n"
	ret = ret + "A=M-D" + "\n"
	ret = ret + "D=M" + "\n"
	ret = ret + "@R3" + "\n"
	ret = ret + "M=D" + "\n"
	// ARG = *(FRAME-3)
	ret = ret + "@3" + "\n"
	ret = ret + "D=A" + "\n"
	ret = ret + "@R5" + "\n"
	ret = ret + "A=M-D" + "\n"
	ret = ret + "D=M" + "\n"
	ret = ret + "@R2" + "\n"
	ret = ret + "M=D" + "\n"
	// LCL = *(FRAME-4)
	ret = ret + "@4" + "\n"
	ret = ret + "D=A" + "\n"
	ret = ret + "@R5" + "\n"
	ret = ret + "A=M-D" + "\n"
	ret = ret + "D=M" + "\n"
	ret = ret + "@R1" + "\n"
	ret = ret + "M=D" + "\n"
	// RET = *(FRAME-5)
	ret = ret + "@5" + "\n"
	ret = ret + "D=A" + "\n"
	ret = ret + "@R5" + "\n"
	ret = ret + "A=M-D" + "\n"
	// goto RET
	ret = ret + "0;JMP" + "\n"
	return ret
}
func writeFunction() string {
	ret := ""
	ret = ret + "@" + arg1() + "\n"
	a2, _ := strconv.Atoi(arg2())
	for i := 0; i < a2; i++ {
		ret = ret + "@0" + "\n"
		ret = ret + "D=A" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "A=M" + "\n"
		ret = ret + "M=D" + "\n"
		ret = ret + "@R0" + "\n"
		ret = ret + "M=M+1" + "\n" // SPの更新
	}
	return ret
}
