package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	aCommand = iota
	cCcmmand
	lCcmmand
)

var commands []string
var pos = 0

var m map[string]int

func main() {
	filename := os.Args[1]

	initialize(filename)
	m = make(map[string]int)
	initializeMap()

	address := 0
	for range commands {
		switch commandType() {
		case aCommand:
			address++
		case cCcmmand:
			address++
		case lCcmmand:
			symbol := symbol()
			m[symbol] = address
		}
		pos++
	}
	pos = 0

	var results string
	varCount := 16
	// for _, c := range commands {
	for range commands {
		// fmt.Println(c)
		if !hasMoreCommands() {
			break
		}
		switch commandType() {
		case aCommand:
			symbol := symbol()
			var b int
			i, err := strconv.Atoi(symbol)
			if err == nil {
				b = i
			} else {
				if val, ok := m[symbol]; ok {
					b = val
				} else {
					b = varCount
					m[symbol] = varCount
					varCount++
				}
			}
			s := fmt.Sprintf("%016b", b)
			results = results + s
		case cCcmmand:
			pre := "111"
			comp := comp()
			dest := dest()
			jump := jump()
			s := pre + comp + dest + jump
			results = results + s
		case lCcmmand:
			pos++
			continue
		}
		results = results + "\n"
		pos++
	}

	filename = strings.Replace(filename, ".asm", ".hack", 1)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("error")
	}
	defer file.Close()

	file.WriteString(results)
}

func initialize(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, " ", "", -1)
		slice := strings.Split(line, "//")
		if slice[0] != "" {
			commands = append(commands, slice[0])
		}
	}
}
func initializeMap() {
	m["SP"] = 0
	m["LCL"] = 1
	m["ARG"] = 2
	m["THIS"] = 3
	m["THAT"] = 4

	m["R0"] = 0
	m["R1"] = 1
	m["R2"] = 2
	m["R3"] = 3
	m["R4"] = 4
	m["R5"] = 5
	m["R6"] = 6
	m["R7"] = 7
	m["R8"] = 8
	m["R9"] = 9
	m["R10"] = 10
	m["R11"] = 11
	m["R12"] = 12
	m["R13"] = 13
	m["R14"] = 14
	m["R15"] = 15

	m["SCREEN"] = 16384
	m["KBD"] = 24576
}

func hasMoreCommands() bool {
	if pos < len(commands) {
		return true
	}
	return false
}

func advance() string {
	if !hasMoreCommands() {
		return ""
	}
	pos++
	return commands[pos]
}

func commandType() int {
	c := commands[pos]
	if strings.Contains(c, "@") {
		return aCommand
	}
	if strings.Contains(c, "=") || strings.Contains(c, ";") {
		return cCcmmand
	}
	return lCcmmand
}

func symbol() string {
	c := commands[pos]
	if strings.Contains(c, "@") {
		return strings.Split(c, "@")[1]
	}
	if strings.Contains(c, "(") {
		c = strings.Replace(c, "(", "", 1)
		c = strings.Replace(c, ")", "", 1)
	}
	return c
}
func dest() string {
	c := commands[pos]
	if !strings.Contains(c, "=") {
		return "000"
	}
	str := strings.Split(c, "=")[0]
	switch str {
	case "M":
		return "001"
	case "D":
		return "010"
	case "MD":
		return "011"
	case "A":
		return "100"
	case "AM":
		return "101"
	case "AD":
		return "110"
	case "AMD":
		return "111"
	}
	return "000"
}
func comp() string {
	c := commands[pos]
	if strings.Contains(c, "=") {
		c = strings.Split(c, "=")[1]
	}
	if strings.Contains(c, ";") {
		c = strings.Split(c, ";")[0]
	}

	ret := "000000"
	switch c {
	case "0":
		ret = "101010"
	case "1":
		ret = "111111"
	case "-1":
		ret = "111010"
	case "D":
		ret = "001100"
	case "A":
		fallthrough
	case "M":
		ret = "110000"
	case "!D":
		ret = "001101"
	case "!A":
		fallthrough
	case "!M":
		ret = "110001"
	case "-D":
		ret = "001111"
	case "-A":
		fallthrough
	case "-M":
		ret = "110011"
	case "D+1":
		ret = "011111"
	case "A+1":
		fallthrough
	case "M+1":
		ret = "110111"
	case "D-1":
		ret = "001110"
	case "A-1":
		fallthrough
	case "M-1":
		ret = "110010"
	case "D+A":
		fallthrough
	case "D+M":
		ret = "000010"
	case "D-A":
		fallthrough
	case "D-M":
		ret = "010011"
	case "A-D":
		fallthrough
	case "M-D":
		ret = "000111"
	case "D&A":
		fallthrough
	case "D&M":
		ret = "000000"
	case "D|A":
		fallthrough
	case "D|M":
		ret = "010101"
	}
	if strings.Contains(c, "M") {
		return "1" + ret
	}
	return "0" + ret
}

func jump() string {
	c := commands[pos]
	if !strings.Contains(c, ";") {
		return "000"
	}
	str := strings.Split(c, ";")[1]
	switch str {
	case "JGT":
		return "001"
	case "JEQ":
		return "010"
	case "JGE":
		return "011"
	case "JLT":
		return "100"
	case "JNE":
		return "101"
	case "JLE":
		return "110"
	case "JMP":
		return "111"
	}
	return "000"
}
