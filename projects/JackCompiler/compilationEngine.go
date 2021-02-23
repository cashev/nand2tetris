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
	cur = skip(cur, CLASS)
	compileIdentifier(cur) // className
	cur = nextToken()
	compileToken(cur) // '{'
	cur = skip(cur, LBRACE)

	for !equal(cur, RBRACE) {
		if equal(cur, STATIC) || equal(cur, FIELD) {
			compileClassVarDec()
		}
		if isSubroutine(cur.str) {
			compileSubroutine()
		}
	}
	compileToken(cur)
	cur = skip(cur, RBRACE)
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
	cur = nextToken()
	for cur.str == COMMA {
		compileToken(cur) // ','
		cur = skip(cur, COMMA)
		compileToken(cur) // varName
		cur = nextToken()
	}
	compileToken(cur) // ';'
	cur = skip(cur, SEMICOLON)
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
	cur = skip(cur, LPAREN)

	compileParameterList() // parameterList

	compileToken(cur) // ')'
	cur = skip(cur, RPAREN)

	// subroutineBody = '{' varDec* statements '}'
	results = append(results, "<subroutineBody>")
	compileToken(cur) // '{'
	cur = skip(cur, LBRACE)
	for equal(cur, VAR) {
		results = append(results, "<varDec>")
		// varDec = 'var' type varName (',' varName)* ';'
		compileToken(cur) // var
		cur = skip(cur, VAR)
		compileToken(cur) // type
		cur = nextToken()
		compileToken((cur)) // varName
		cur = nextToken()
		for equal(cur, COMMA) {
			compileToken(cur) // ','
			cur = skip(cur, COMMA)
			compileToken((cur)) // varName
			cur = nextToken()
		}
		compileToken(cur) // ';'
		cur = skip(cur, SEMICOLON)
		results = append(results, "</varDec>")
	}
	cur = compileStatements()

	compileToken(cur) // '}'
	cur = skip(cur, RBRACE)
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
		cur = skip(cur, COMMA)
		compileToken(cur) // type
		cur = nextToken()
		compileToken(cur) // varName
		cur = nextToken()
	}
	return
}

func compileVarDec() {

}

// statements = statement*
func compileStatements() token {
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
	return cur
}

// doStatement = 'do' subroutineCall ';'
func compileDo() {
	results = append(results, "<doStatement>")
	compileToken(cur) // 'do'
	cur = skip(cur, DO)
	compileToken(cur) // subroutineName | (className | varName)
	cur = nextToken()
	if cur.str == LPAREN {
		// subroutineCall = subroutineName '(' exprList ')'
		compileToken(cur) // '('
		cur = skip(cur, LPAREN)
		cur = compileExpressionList()
		compileToken(cur) // ')'
		cur = skip(cur, RPAREN)
	}
	if cur.str == PERIOD {
		// subroutineCall = (className | varName) '.' subroutineName
		// 									'(' exprList ')'
		compileToken(cur) // '.'
		cur = skip(cur, PERIOD)
		compileToken(cur) // identifier
		cur = nextToken()
		compileToken(cur) // '('
		cur = nextToken()
		cur = compileExpressionList()
		compileToken(cur)
		cur = skip(cur, RPAREN) // ')'
	}
	compileToken(cur) // ';'
	cur = skip(cur, SEMICOLON)
	results = append(results, "</doStatement>")
}

// letstatement = 'let' varName ('[' expr ']')? '=' expr ';'
func compileLet() {
	results = append(results, "<letStatement>")
	compileToken(cur) // let
	cur = skip(cur, LET)
	compileToken(cur) // varName
	cur = nextToken()
	if equal(cur, LBRACKET) {
		compileToken(cur) // '['
		cur = skip(cur, LBRACKET)
		compileExpression()
		compileToken(cur) // ']'
		cur = skip(cur, RBRACKET)
	}
	compileToken(cur) // '='
	cur = skip(cur, EQUAL)
	compileExpression()
	compileToken(cur) // ';'
	cur = skip(cur, SEMICOLON)
	results = append(results, "</letStatement>")
}

// whileStatement = 'while' '(' expr ')' '{' statements '}'
func compileWhile() {
	results = append(results, "<whileStatement>")
	compileToken(cur) // 'while'
	cur = skip(cur, WHILE)
	compileToken(cur) // '('
	cur = skip(cur, LPAREN)
	compileExpression()
	compileToken(cur) // ')'
	cur = skip(cur, RPAREN)
	compileToken(cur) // '{'
	cur = skip(cur, LBRACE)
	compileStatements()
	compileToken(cur) // '}'
	cur = skip(cur, RBRACE)
	results = append(results, "</whileStatement>")
}

// returnStatement = 'return' expr? ';'
func compileReturn() {
	results = append(results, "<returnStatement>")
	compileToken(cur) // 'return'
	cur = skip(cur, RETURN)
	if !equal(cur, SEMICOLON) {
		compileExpression()
	}
	compileToken(cur) // ';'
	cur = skip(cur, SEMICOLON)
	results = append(results, "</returnStatement>")
}

// ifStatement = 'if' '(' expr ')' '{' statements '}'
// 							 ('else' '{' statements '}')?
func compileIf() {
	results = append(results, "<ifStatement>")
	compileToken(cur) // 'if'
	cur = skip(cur, IF)
	compileToken(cur) // '('
	cur = skip(cur, LPAREN)
	compileExpression()
	compileToken(cur) // ')'
	cur = skip(cur, RPAREN)
	compileToken(cur) // '{'
	cur = skip(cur, LBRACE)
	compileStatements()
	compileToken(cur) // '}'
	cur = skip(cur, RBRACE)
	if equal(cur, ELSE) {
		compileToken(cur) // 'else'
		cur = skip(cur, ELSE)
		compileToken(cur) // '{'
		cur = skip(cur, LBRACE)
		compileStatements()
		compileToken(cur) // '}'
		cur = skip(cur, RBRACE)
	}
	results = append(results, "</ifStatement>")
}

// expr = term (op term)*
func compileExpression() {
	results = append(results, "<expression>")
	compileTerm() // term
	cur = nextToken()
	for isOperator(cur.str) {
		compileToken(cur) // op
		cur = nextToken()
		compileTerm() // term
		cur = nextToken()
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
	now := cur
	results = append(results, "<term>")
	if now.kind == IDENTIFIER {
		compileToken(now)
		if readNextToken().str == LBRACKET {
			// varName '[' expr ']'
			compileToken(nextToken()) // '['
			cur = nextToken()
			compileExpression()
			compileToken(cur) // ']'
		}
		if readNextToken().str == LPAREN {
			// subroutineCall = subroutineName '(' exprList ')'
			compileToken(nextToken()) // '('
			cur = nextToken()
			next := compileExpressionList()
			compileToken(next) // ')'

		}
		if readNextToken().str == PERIOD {
			// subroutineCall = (className | varName) '.' subroutineName
			// 									'(' exprList ')'
			compileToken(nextToken()) // '.'
			compileToken(nextToken()) // subroutineName
			compileToken(nextToken()) // '('
			cur = nextToken()
			next := compileExpressionList()
			compileToken(next) // ')'
		}

	} else if isUnaryOperator(now.str) {
		compileToken(now) // uparyOp
		cur = nextToken()
		compileTerm() // term
	} else if now.str == LPAREN {
		// '(' expr ')'
		compileToken(now) // '('
		cur = nextToken()
		compileExpression()
		compileToken(cur) // ')'
	} else {
		compileToken(now)
	}
	results = append(results, "</term>")
}

// exprList = (expr (',' expr)* )?
func compileExpressionList() token {
	results = append(results, "<expressionList>")
	now := cur
	if now.str == RPAREN {
		results = append(results, "</expressionList>")
		return now
	}
	cur = now
	compileExpression()

	tok := cur
	for tok.str == COMMA {
		compileToken(tok) // ','
		cur = nextToken()
		compileExpression()
		tok = cur
	}
	results = append(results, "</expressionList>")
	return tok
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

func readNextToken() token {
	return tokens[position+1]
}

func equal(tok token, str string) bool {
	return tok.str == str
}

func skip(tok token, str string) token {
	if !equal(tok, str) {
		msg := fmt.Sprintf("skip error. unexpected token. expected %s", str)
		panic(msg)
	}
	position++
	return tokens[position]
}
