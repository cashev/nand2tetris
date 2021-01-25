package main

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
