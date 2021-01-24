package main

import (
	"fmt"
	"strings"
)

var input = ""
var pos = 0

type tokenKind int

const (
	KEYWORD tokenKind = iota + 1
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

const (
	// キーワード
	CLASS       = "class"
	CONSTRUCTOR = "constructor"
	FUNCTION    = "function"
	METHOD      = "method"
	FIELD       = "field"
	STATIC      = "static"
	VAR         = "var"
	INT         = "int"
	CHAR        = "char"
	BOOLEAN     = "boolean"
	VOID        = "void"
	TRUE        = "true"
	FALSE       = "false"
	NULL        = "null"
	THIS        = "this"
	LET         = "let"
	DO          = "do"
	IF          = "if"
	ELSE        = "else"
	WHILE       = "while"
	RETURN      = "return"

	// 記号
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	PERIOD    = "."
	COMMA     = ","
	SEMICOLON = ";"
	PLUS      = "+"
	MINUS     = "-"
	ASTERISK  = "*"
	SLASH     = "/"
	AND       = "&"
	OR        = "|"
	LT        = "<"
	RT        = ">"
	EQUAL     = "="
	TILDE     = "~"

	COMMENT         = "//"
	COMMENTSTART    = "/*"
	APICOMMENTSTART = "/**"
	COMMENTEND      = "*/"
)

type token struct {
	kind tokenKind
	str  string
}

func new(kind tokenKind, str string) token {
	return token{kind: kind, str: str}
}

func initialize(in []string) {
	in = removeComment(in)

	pos = 0
	for _, str := range in {
		input += str
	}
}

func tokenize() []token {
	var tokens []token

	return tokens
}

func hasMoreTokens() bool {
	return pos < len(input)
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
		fmt.Println(line)
		str := strings.Split(line, COMMENT)[0]
		if str != "" {
			results = append(results, str)
		}
	}

	return results
}
