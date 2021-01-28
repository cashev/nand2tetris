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

	for hasMoreTokens() {
		now := input[pos]
		if isSpace(now) {
			pos++
			continue
		}

		switch now {
		case '+':
			tok := new(SYMBOL, PLUS)
			tokens = append(tokens, tok)
		case '-':
			tok := new(SYMBOL, MINUS)
			tokens = append(tokens, tok)
		case '*':
			tok := new(SYMBOL, ASTERISK)
			tokens = append(tokens, tok)
		case '/':
			tok := new(SYMBOL, SLASH)
			tokens = append(tokens, tok)
		case '&':
			tok := new(SYMBOL, AND)
			tokens = append(tokens, tok)
		case '|':
			tok := new(SYMBOL, OR)
			tokens = append(tokens, tok)
		case '<':
			tok := new(SYMBOL, LT)
			tokens = append(tokens, tok)
		case '>':
			tok := new(SYMBOL, RT)
			tokens = append(tokens, tok)
		case '=':
			tok := new(SYMBOL, EQUAL)
			tokens = append(tokens, tok)
		case '.':
			tok := new(SYMBOL, PERIOD)
			tokens = append(tokens, tok)
		case ',':
			tok := new(SYMBOL, COMMA)
			tokens = append(tokens, tok)
		case ';':
			tok := new(SYMBOL, SEMICOLON)
			tokens = append(tokens, tok)
		case '(':
			tok := new(SYMBOL, LPAREN)
			tokens = append(tokens, tok)
		case ')':
			tok := new(SYMBOL, RPAREN)
			tokens = append(tokens, tok)
		case '{':
			tok := new(SYMBOL, LBRACE)
			tokens = append(tokens, tok)
		case '}':
			tok := new(SYMBOL, RBRACE)
			tokens = append(tokens, tok)
		case '"':
			now = next()
			str := string(now)
			for readNext() != '"' {
				str = str + string(next())
			}
			tok := new(STRING_CONST, str)
			tokens = append(tokens, tok)
			next() // '"'を読み飛ばす
		default:
			if isDigit(now) {
				str := string(now)
				for isDigit(readNext()) {
					str = str + string(next())
				}
				tok := new(INT_CONST, str)
				tokens = append(tokens, tok)
			} else if isLetter(now) {
				str := string(now)
				for isLetter(readNext()) || isDigit(readNext()) {
					str = str + string(next())
				}
				if isKeyword(str) {
					tok := new(KEYWORD, str)
					tokens = append(tokens, tok)
					pos++
					continue
				}
				tok := new(IDENTIFIER, str)
				tokens = append(tokens, tok)
			}
		}
		pos++
	}

	return tokens
}

func hasMoreTokens() bool {
	return pos < len(input)
}

func next() byte {
	if !hasMoreTokens() {
		return ' '
	}
	pos++
	return input[pos]
}

func readNext() byte {
	if !hasMoreTokens() {
		return ' '
	}
	return input[pos+1]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isKeyword(str string) bool {
	keywords := []string{CLASS, CONSTRUCTOR, FUNCTION, METHOD, FIELD, STATIC, VAR, INT, CHAR,
		BOOLEAN, VOID, TRUE, FALSE, NULL, THIS, LET, DO, IF, ELSE, WHILE, RETURN}
	for _, keyword := range keywords {
		if keyword == str {
			return true
		}
	}
	return false
}

func isSpace(ch byte) bool {
	return ch == ' '
}
