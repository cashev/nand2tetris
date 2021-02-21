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
	cur = compileStatements(cur)

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
func compileStatements(now token) token {
	results = append(results, "<statements>")
	// statement = letStatement | ifStatement |
	// 						 whileStatement | doStatement |
	// 						 returnStatement
	loop := true
	tok := now
	for {
		switch tok.str {
		case LET:
			compileLet(tok)
		case IF:
			compileIf(tok)
		case WHILE:
			compileWhile(tok)
		case DO:
			compileDo(tok)
		case RETURN:
			compileReturn(tok)
		default:
			loop = false
		}
		if !loop {
			break
		}
		tok = nextToken()
	}
	results = append(results, "</statements>")
	return tok
}

// doStatement = 'do' subroutineCall ';'
func compileDo(now token) {
	results = append(results, "<doStatement>")
	compileToken(now)         // do
	compileToken(nextToken()) // identifier
	tok := nextToken()
	if tok.str == LPAREN {
		// subroutineCall = subroutineName '(' exprList ')'
		compileToken(tok) // '('
		next := compileExpressionList(nextToken())
		compileToken(next) // ')'
	}
	if tok.str == PERIOD {
		// subroutineCall = (className | varName) '.' subroutineName
		// 									'(' exprList ')'
		compileToken(tok)         // '.'
		compileToken(nextToken()) // identifier
		compileToken(nextToken()) // '('
		next := compileExpressionList(nextToken())
		compileToken(next) // ')'
	}
	compileToken(nextToken()) // ';'

	results = append(results, "</doStatement>")
}

// letstatement = 'let' varName ('[' expr ']')? '=' expr ';'
func compileLet(now token) {
	results = append(results, "<letStatement>")
	compileToken(now)         // let
	compileToken(nextToken()) // varName
	tok := nextToken()
	if tok.str == LBRACKET {
		compileToken(tok) // '['
		tok1 := compileExpression(nextToken())
		compileToken(tok1) // ']'
		tok = nextToken()
	}
	compileToken(tok) // '='
	tok1 := compileExpression(nextToken())
	compileToken(tok1) // ';'

	results = append(results, "</letStatement>")
}

// whileStatement = 'while' '(' expr ')' '{' statements '}'
func compileWhile(now token) {
	results = append(results, "<whileStatement>")
	compileToken(now)         // 'while'
	compileToken(nextToken()) // '('
	tok := compileExpression(nextToken())
	compileToken(tok)         // ')'
	compileToken(nextToken()) // '{'
	tok = compileStatements(nextToken())
	compileToken(tok) // '}'
	results = append(results, "</whileStatement>")
}

// returnStatement = 'return' expr? ';'
func compileReturn(now token) {
	results = append(results, "<returnStatement>")
	compileToken(now) // 'return'
	tok := nextToken()
	if tok.str != SEMICOLON {
		tok = compileExpression(tok)
	}
	compileToken(tok) // ';'
	results = append(results, "</returnStatement>")
}

// ifStatement = 'if' '(' expr ')' '{' statements '}'
// 							 ('else' '{' statements '}')?
func compileIf(now token) {
	results = append(results, "<ifStatement>")
	compileToken(now)         // 'if'
	compileToken(nextToken()) // '('
	tok := compileExpression(nextToken())
	compileToken(tok)         // ')'
	compileToken(nextToken()) // '{'
	tok = compileStatements(nextToken())
	compileToken(tok) // '}'
	if readNextToken().str == ELSE {
		compileToken(nextToken()) // 'else'
		compileToken(nextToken()) // '{'
		tok = compileStatements(nextToken())
		compileToken(tok) // '}'
	}
	results = append(results, "</ifStatement>")
}

// expr = term (op term)*
func compileExpression(now token) token {
	results = append(results, "<expression>")
	compileTerm(now) // term
	tok := nextToken()
	for isOp(tok.str) {
		compileToken(tok)        // op
		compileTerm(nextToken()) // term
		tok = nextToken()
	}
	results = append(results, "</expression>")
	return tok
}

func isOp(str string) bool {
	ops := []string{PLUS, MINUS, ASTERISK, SLASH, AND, OR, LT, RT, EQUAL}
	for _, op := range ops {
		if op == str {
			return true
		}
	}
	return false
}

func isUnaryOp(str string) bool {
	return str == "-" || str == "~"
}

// term = integerConstant | stringConstant | keywordConstant
// 				| varName | varName '[' expr ']' | subroutineCall
// 				| '(' expr ')' | unaryOp term
func compileTerm(now token) {
	results = append(results, "<term>")
	if now.kind == IDENTIFIER {
		compileToken(now)
		if readNextToken().str == LBRACKET {
			// varName '[' expr ']'
			compileToken(nextToken()) // '['
			tok := compileExpression(nextToken())
			compileToken(tok) // ']'
		}
		if readNextToken().str == LPAREN {
			// subroutineCall = subroutineName '(' exprList ')'
			compileToken(nextToken()) // '('
			next := compileExpressionList(nextToken())
			compileToken(next) // ')'

		}
		if readNextToken().str == PERIOD {
			// subroutineCall = (className | varName) '.' subroutineName
			// 									'(' exprList ')'
			compileToken(nextToken()) // '.'
			compileToken(nextToken()) // subroutineName
			compileToken(nextToken()) // '('
			next := compileExpressionList(nextToken())
			compileToken(next) // ')'
		}

	} else if isUnaryOp(now.str) {
		compileToken(now)        // uparyOp
		compileTerm(nextToken()) // term
	} else if now.str == LPAREN {
		// '(' expr ')'
		compileToken(now) // '('
		tok := compileExpression(nextToken())
		compileToken(tok) // ')'
	} else {
		compileToken(now)
	}
	results = append(results, "</term>")
}

// exprList = (expr (',' expr)* )?
func compileExpressionList(now token) token {
	results = append(results, "<expressionList>")
	if now.str == RPAREN {
		results = append(results, "</expressionList>")
		return now
	}
	tok1 := compileExpression(now)

	tok := tok1
	for tok.str == COMMA {
		compileToken(tok) // ','
		tok2 := compileExpression(nextToken())
		tok = tok2
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
