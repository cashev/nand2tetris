package main

import (
	"fmt"
	"strconv"
)

var position int
var tokens []token
var cur token

var results []string

type identifierKind string

const (
	iStatic identifierKind = "STATIC"
	iField  identifierKind = "FIELD"
	iArg    identifierKind = "ARG"
	iVar    identifierKind = "VAR"
	iNone   identifierKind = "NONE"
)

type identifier struct {
	name string
	typ  string
	kind identifierKind
}

var symbolTable []identifier
var subroutineSymbolTable []identifier

const (
	EOF = "END OF FILE"

	DEFINED = "defined"
	USED    = "used"

	ARGUMENT   = "argument"
	SUBROUTINE = "subroutine"
)

func initializeCompileEngine() {
	position = 0
	tokens = make([]token, 0)
	cur = token{}
	results = make([]string, 0)

	symbolTable = make([]identifier, 0)
	subroutineSymbolTable = make([]identifier, 0)
}

func compile(toks []token) []string {
	initializeCompileEngine()
	tokens = append(toks, token{str: EOF})
	cur = tokens[0]

	for !equal(cur, EOF) {
		compileClass()
	}

	return results
}

// class = 'class' className '{' classVarDec* subroutineDec* '}'
func compileClass() {
	symbolTable = make([]identifier, 0)

	results = append(results, "<class>")
	compileToken(cur) // class
	consume(CLASS)
	compileIdentifier(cur, CLASS, DEFINED) // className
	consumeIdentifier()
	compileToken(cur) // '{'
	consume(LBRACE)

	for !equal(cur, RBRACE) {
		if equal(cur, STATIC) || equal(cur, FIELD) {
			compileClassVarDec()
		}
		if isSubroutine(cur.str) {
			compileSubroutine()
		}
	}
	compileToken(cur)
	consume(RBRACE)
	results = append(results, "</class>")
}

// classVarDec = ('static' | 'field') type varName (',' varName)* ';'
func compileClassVarDec() {
	results = append(results, "<classVarDec>")
	compileToken(cur) // 'static' | 'field'
	category := cur.str
	cur = nextToken()
	compileToken(cur) // type
	typ := cur.str
	cur = nextToken()
	// varName
	var iKind identifierKind
	if category == STATIC {
		iKind = iStatic
	} else {
		iKind = iField
	}
	define(cur.str, typ, iKind)
	compileIdentifier(cur, category, DEFINED)
	consumeIdentifier()
	for cur.str == COMMA {
		compileToken(cur) // ','
		consume(COMMA)
		// varName
		define(cur.str, typ, iKind)
		compileIdentifier(cur, category, DEFINED)
		cur = nextToken()
	}
	compileToken(cur) // ';'
	consume(SEMICOLON)
	results = append(results, "</classVarDec>")
}

func isSubroutine(str string) bool {
	return str == FUNCTION || str == METHOD || str == CONSTRUCTOR
}

// subroutineDec = ('construct' | 'function' | 'method')
//									('void' | type) subroutineName '(' parameterList ')'
//									subroutineBody
func compileSubroutine() {
	subroutineSymbolTable = make([]identifier, 0)

	results = append(results, "<subroutineDec>")
	compileToken(cur) // 'construct' | 'function' | 'method'
	cur = nextToken()
	compileToken(cur) // 'void' | type
	cur = nextToken()
	// subroutineName
	compileIdentifier(cur, SUBROUTINE, DEFINED)
	cur = nextToken()
	compileToken(cur) // '('
	consume(LPAREN)

	compileParameterList() // parameterList

	compileToken(cur) // ')'
	consume(RPAREN)

	// subroutineBody = '{' varDec* statements '}'
	results = append(results, "<subroutineBody>")
	compileToken(cur) // '{'
	consume(LBRACE)
	for equal(cur, VAR) {
		results = append(results, "<varDec>")
		// varDec = 'var' type varName (',' varName)* ';'
		compileToken(cur) // 'var'
		consume(VAR)
		compileToken(cur) // type
		typ := cur.str
		cur = nextToken()
		// varName
		define(cur.str, typ, iVar)
		compileIdentifier(cur, VAR, DEFINED)
		cur = nextToken()
		for equal(cur, COMMA) {
			compileToken(cur) // ','
			consume(COMMA)
			// varName
			define(cur.str, typ, iVar)
			compileIdentifier(cur, VAR, DEFINED)
			cur = nextToken()
		}
		compileToken(cur) // ';'
		consume(SEMICOLON)
		results = append(results, "</varDec>")
	}
	compileStatements()

	compileToken(cur) // '}'
	consume(RBRACE)
	results = append(results, "</subroutineBody>")
	results = append(results, "</subroutineDec>")
}

// parameterList = ((type varName) (',' type varName)* )?
func compileParameterList() {
	results = append(results, "<parameterList>")
	if equal(cur, RPAREN) {
		results = append(results, "</parameterList>")
		return
	}
	compileToken(cur) // type
	typ := cur.str
	cur = nextToken()
	// varName
	define(cur.str, typ, iArg)
	compileIdentifier(cur, ARGUMENT, DEFINED)
	cur = nextToken()
	for equal(cur, COMMA) {
		compileToken(cur) // ','
		consume(COMMA)
		compileToken(cur) // type
		cur = nextToken()
		compileToken(cur) // varName
		cur = nextToken()
	}
}

// statements = statement*
func compileStatements() {
	results = append(results, "<statements>")
	// statement = letStatement | ifStatement |
	// 						 whileStatement | doStatement |
	// 						 returnStatement
	loop := true
	for loop {
		switch cur.str {
		case LET:
			compileLet()
		case IF:
			compileIf()
		case WHILE:
			compileWhile()
		case DO:
			compileDo()
		case RETURN:
			compileReturn()
		default:
			loop = false
		}
	}
	results = append(results, "</statements>")
}

// doStatement = 'do' subroutineCall ';'
func compileDo() {
	results = append(results, "<doStatement>")
	compileToken(cur) // 'do'
	consume(DO)
	next := readNextToken()
	if equal(next, LPAREN) {
		// subroutineCall = subroutineName '(' exprList ')'
		// subroutineName
		compileIdentifier(cur, SUBROUTINE, USED)
		cur = nextToken()
		// '('
		compileToken(cur)
		consume(LPAREN)
		compileExpressionList()
		compileToken(cur) // ')'
		consume(RPAREN)
	} else if equal(next, PERIOD) {
		// subroutineCall = (className | varName) '.' subroutineName
		// 									'(' exprList ')'
		// className | varName
		kind := kindOf(cur.str)
		var category string
		if kind == iNone {
			category = CLASS
		} else {
			category = string(kind)
		}
		compileIdentifier(cur, category, USED)
		cur = nextToken()
		// '.'
		compileToken(cur)
		consume(PERIOD)
		// subroutineName
		compileIdentifier(cur, SUBROUTINE, USED)
		cur = nextToken()
		compileToken(cur) // '('
		cur = nextToken()
		compileExpressionList()
		compileToken(cur) // ')'
		consume(RPAREN)
	}
	compileToken(cur) // ';'
	consume(SEMICOLON)
	results = append(results, "</doStatement>")
}

// letstatement = 'let' varName ('[' expr ']')? '=' expr ';'
func compileLet() {
	results = append(results, "<letStatement>")
	compileToken(cur) // let
	consume(LET)
	// varName
	compileIdentifier(cur, string(kindOf(cur.str)), USED)
	cur = nextToken()
	if equal(cur, LBRACKET) {
		compileToken(cur) // '['
		consume(LBRACKET)
		compileExpression()
		compileToken(cur) // ']'
		consume(RBRACKET)
	}
	compileToken(cur) // '='
	consume(EQUAL)
	compileExpression()
	compileToken(cur) // ';'
	consume(SEMICOLON)
	results = append(results, "</letStatement>")
}

// whileStatement = 'while' '(' expr ')' '{' statements '}'
func compileWhile() {
	results = append(results, "<whileStatement>")
	compileToken(cur) // 'while'
	consume(WHILE)
	compileToken(cur) // '('
	consume(LPAREN)
	compileExpression()
	compileToken(cur) // ')'
	consume(RPAREN)
	compileToken(cur) // '{'
	consume(LBRACE)
	compileStatements()
	compileToken(cur) // '}'
	consume(RBRACE)
	results = append(results, "</whileStatement>")
}

// returnStatement = 'return' expr? ';'
func compileReturn() {
	results = append(results, "<returnStatement>")
	compileToken(cur) // 'return'
	consume(RETURN)
	if !equal(cur, SEMICOLON) {
		compileExpression()
	}
	compileToken(cur) // ';'
	consume(SEMICOLON)
	results = append(results, "</returnStatement>")
}

// ifStatement = 'if' '(' expr ')' '{' statements '}'
// 							 ('else' '{' statements '}')?
func compileIf() {
	results = append(results, "<ifStatement>")
	compileToken(cur) // 'if'
	consume(IF)
	compileToken(cur) // '('
	consume(LPAREN)
	compileExpression()
	compileToken(cur) // ')'
	consume(RPAREN)
	compileToken(cur) // '{'
	consume(LBRACE)
	compileStatements()
	compileToken(cur) // '}'
	consume(RBRACE)
	if equal(cur, ELSE) {
		compileToken(cur) // 'else'
		consume(ELSE)
		compileToken(cur) // '{'
		consume(LBRACE)
		compileStatements()
		compileToken(cur) // '}'
		consume(RBRACE)
	}
	results = append(results, "</ifStatement>")
}

// expr = term (op term)*
func compileExpression() {
	results = append(results, "<expression>")
	compileTerm() // term
	for isOperator(cur.str) {
		compileToken(cur) // op
		cur = nextToken()
		compileTerm() // term
	}
	results = append(results, "</expression>")
}

func isOperator(str string) bool {
	ops := []string{PLUS, MINUS, ASTERISK, SLASH, AND, OR, LT, RT, EQUAL}
	for _, op := range ops {
		if op == str {
			return true
		}
	}
	return false
}

func isUnaryOperator(str string) bool {
	return str == "-" || str == "~"
}

// term = integerConstant | stringConstant | keywordConstant
// 				| varName | varName '[' expr ']' | subroutineCall
// 				| '(' expr ')' | unaryOp term
func compileTerm() {
	results = append(results, "<term>")
	if cur.kind == IDENTIFIER {
		next := readNextToken()
		if equal(next, LBRACKET) {
			// varName | varName '[' expr ']'
			// varName
			compileIdentifier(cur, string(kindOf(cur.str)), USED)
			cur = nextToken()
			compileToken(cur) // '['
			consume(LBRACKET)
			compileExpression()
			compileToken(cur) // ']'
			consume(RBRACKET)
		} else if equal(next, LPAREN) {
			// subroutineCall = subroutineName '(' exprList ')'
			// subroutineName
			compileIdentifier(cur, string(kindOf(cur.str)), USED)
			cur = nextToken()
			// '('
			compileToken(cur)
			consume(LPAREN)
			compileExpressionList()
			compileToken(cur) // ')'
			consume(RPAREN)
		} else if equal(next, PERIOD) {
			// subroutineCall = (className | varName) '.' subroutineName
			// 									'(' exprList ')'
			// className | varName
			kind := kindOf(cur.str)
			var category string
			if kind == iNone {
				category = CLASS
			} else {
				category = string(kind)
			}
			compileIdentifier(cur, category, USED)
			cur = nextToken()
			// '.'
			compileToken(cur)
			consume(PERIOD)
			// subroutineName
			compileIdentifier(cur, SUBROUTINE, USED)
			cur = nextToken()
			compileToken(cur) // '('
			consume(LPAREN)
			compileExpressionList()
			compileToken(cur) // ')'
			consume(RPAREN)
		} else {
			compileIdentifier(cur, string(kindOf(cur.str)), USED)
			cur = nextToken()
		}
	} else if isUnaryOperator(cur.str) {
		compileToken(cur) // unaryOperator
		cur = nextToken()
		compileTerm() // term
	} else if equal(cur, LPAREN) {
		// '(' expr ')'
		compileToken(cur) // '('
		consume(LPAREN)
		compileExpression()
		compileToken(cur) // ')'
		consume(RPAREN)
	} else {
		// integerConstant | stringConstant | keywordConsant
		compileToken(cur)
		cur = nextToken()
	}
	results = append(results, "</term>")
}

// exprList = (expr (',' expr)* )?
func compileExpressionList() {
	results = append(results, "<expressionList>")
	if equal(cur, RPAREN) {
		results = append(results, "</expressionList>")
		return
	}
	compileExpression()
	for equal(cur, COMMA) {
		compileToken(cur) // '.'
		consume(COMMA)
		compileExpression()
	}
	results = append(results, "</expressionList>")
}

func compileToken(tok token) {
	kind := string(tok.kind)
	str := "<" + kind + "> " + tok.str + " </" + kind + ">"
	results = append(results, str)
}

func compileIdentifier(tok token, category string, defOrUsed string) {
	if tok.kind != IDENTIFIER {
		panic("compile error. not identifier.")
	}
	tokKind := string(tok.kind)
	// str := "<" + tokKind + "> " + tok.str + " </" + tokKind + ">"
	str := "<" + tokKind + "> "
	str = str + "name: " + tok.str
	str = str + ", category: " + category
	str = str + ", D/U: " + defOrUsed
	kind := string(kindOf(tok.str))
	str = str + ", kind: " + kind
	index := indexOf(tok.str)
	str = str + ", index: " + strconv.Itoa(index)
	str = str + " </" + tokKind + ">"
	results = append(results, str)
}

func define(name string, t string, kind identifierKind) {
	i := identifier{name: name, typ: t, kind: kind}
	if kind == iStatic || kind == iField {
		symbolTable = append(symbolTable, i)
	}
	if kind == iArg || kind == iVar {
		subroutineSymbolTable = append(subroutineSymbolTable, i)
	}
}

func indexOf(name string) int {
	for i, s := range subroutineSymbolTable {
		if s.name == name {
			return i
		}
	}
	for i, s := range symbolTable {
		if s.name == name {
			return i
		}
	}
	return -1
}

func kindOf(name string) identifierKind {
	for _, s := range subroutineSymbolTable {
		if s.name == name {
			return s.kind
		}
	}
	for _, s := range symbolTable {
		if s.name == name {
			return s.kind
		}
	}
	return iNone
}

func nextToken() token {
	position++
	return tokens[position]
}

func readNextToken() token {
	return tokens[position+1]
}

func equal(tok token, str string) bool {
	return tok.str == str
}

func consume(str string) {
	if !equal(cur, str) {
		msg := fmt.Sprintf("consume error. unexpected token. expected %s", str)
		panic(msg)
	}
	position++
	cur = tokens[position]
}

func consumeIdentifier() {
	if cur.kind != IDENTIFIER {
		msg := fmt.Sprintf("consume error. not identifier token. got %s", cur.str)
		panic(msg)
	}
	position++
	cur = tokens[position]
}
