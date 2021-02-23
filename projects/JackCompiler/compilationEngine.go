package main

import "fmt"

var position int
var tokens []token
var cur token

var results []string

const (
	EOF = "END OF FILE"
)

func initializeCompileEngine() {
	position = 0
	tokens = make([]token, 0)
	cur = token{}
	results = make([]string, 0)
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
	results = append(results, "<class>")
	compileToken(cur) // class
	consume(CLASS)
	compileIdentifier(cur) // className
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
	cur = nextToken()
	compileToken(cur) // type
	cur = nextToken()
	compileToken(cur) // varName
	consumeIdentifier()
	for cur.str == COMMA {
		compileToken(cur) // ','
		consume(COMMA)
		compileToken(cur) // varName
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
	results = append(results, "<subroutineDec>")
	compileToken(cur) // subroutineDec
	cur = nextToken()
	compileToken(cur) // type
	cur = nextToken()
	compileToken(cur) // subroutineName
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
		compileToken(cur) // var
		consume(VAR)
		compileToken(cur) // type
		cur = nextToken()
		compileToken((cur)) // varName
		cur = nextToken()
		for equal(cur, COMMA) {
			compileToken(cur) // ','
			consume(COMMA)
			compileToken((cur)) // varName
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
	cur = nextToken()
	compileToken(cur) // varName
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
	compileToken(cur) // subroutineName | (className | varName)
	cur = nextToken()
	if equal(cur, LPAREN) {
		// subroutineCall = subroutineName '(' exprList ')'
		compileToken(cur) // '('
		consume(LPAREN)
		compileExpressionList()
		compileToken(cur) // ')'
		consume(RPAREN)
	}
	if equal(cur, PERIOD) {
		// subroutineCall = (className | varName) '.' subroutineName
		// 									'(' exprList ')'
		compileToken(cur) // '.'
		consume(PERIOD)
		compileToken(cur) // identifier
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
	compileToken(cur) // varName
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
		compileToken(cur)
		cur = nextToken()
		if equal(cur, LBRACKET) {
			// varName | varName '[' expr ']'
			compileToken(cur) // '['
			consume(LBRACKET)
			compileExpression()
			compileToken(cur) // ']'
			consume(RBRACKET)
		}
		if equal(cur, LPAREN) {
			// subroutineCall = subroutineName '(' exprList ')'
			compileToken(cur) // '('
			consume(LPAREN)
			compileExpressionList()
			compileToken(cur) // ')'
			consume(RPAREN)
		}
		if equal(cur, PERIOD) {
			// subroutineCall = (className | varName) '.' subroutineName
			// 									'(' exprList ')'
			compileToken(cur) // '.'
			consume(PERIOD)
			compileToken(cur) // subroutineName
			cur = nextToken()
			compileToken(cur) // '('
			consume(LPAREN)
			compileExpressionList()
			compileToken(cur) // ')'
			consume(RPAREN)
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

func compileIdentifier(tok token) {
	if tok.kind != IDENTIFIER {
		panic("compile error. not identifier.")
	}
	kind := string(tok.kind)
	str := "<" + kind + "> " + tok.str + " </" + kind + ">"
	results = append(results, str)
}

func nextToken() token {
	position++
	return tokens[position]
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
