package main

var position int
var tokens []token

var results []string

func initializeCompileEngine() {
	position = 0
	tokens = make([]token, 0)
	results = make([]string, 0)
}

func compile(toks []token) []string {
	initializeCompileEngine()
	tokens = toks

	for position+1 < len(tokens) {
		tok := tokens[position]
		compileClass(tok)
	}
	return results
}

// class = 'class' className '{' classVarDec* subroutineDec* '}'
func compileClass(now token) {
	results = append(results, "<class>")
	compileToken(now)         // class
	compileToken(nextToken()) // className
	compileToken(nextToken()) // '{'
	tok := nextToken()
	for tok.str != RBRACE {
		if tok.str == "static" || tok.str == "field" {
			compileClassVarDec(tok)
		}
		if isSubroutine(tok.str) {
			compileSubroutine(tok)
		}
		tok = nextToken()
	}
	compileToken(tok) // '}'
	results = append(results, "</class>")
}

// classVarDec = ('static' | 'field') type varName (',' varName)* ';'
func compileClassVarDec(now token) {
	results = append(results, "<classVarDec>")
	compileToken(now)         // ('static' | 'field')
	compileToken(nextToken()) // type
	compileToken(nextToken()) // varName
	tok := nextToken()
	for tok.str == COMMA {
		compileToken(tok)         // ','
		compileToken(nextToken()) // varName
		tok = nextToken()
	}
	compileToken(tok) // ';'
	results = append(results, "</classVarDec>")
}

func isSubroutine(str string) bool {
	return str == FUNCTION || str == METHOD || str == CONSTRUCTOR
}

// subroutineDec = ('construct' | 'function' | 'method')
//									('void' | type) subroutineName '(' parameterList ')'
//									subroutineBody
func compileSubroutine(now token) {
	results = append(results, "<subroutineDec>")
	compileToken(now)         // subroutineDec
	compileToken(nextToken()) // type
	compileToken(nextToken()) // subroutineName

	compileToken(nextToken()) // '('
	tok := compileParameterList(nextToken())
	compileToken(tok) // ')'

	// subroutineBody = '{' varDec* statements '}'
	results = append(results, "<subroutineBody>")
	compileToken(nextToken()) // '{'
	tok1 := nextToken()
	for tok1.str == "var" {
		results = append(results, "<varDec>")
		// varDec = 'var' type varName (',' varName)* ';'
		compileToken(tok1)        // var
		compileToken(nextToken()) // type
		compileToken(nextToken()) // varName
		tok2 := nextToken()
		for tok2.str == COMMA {
			compileToken(tok2)        // ','
			compileToken(nextToken()) // varName
			tok2 = nextToken()
		}
		compileToken(tok2) // ';'
		tok1 = nextToken()
		results = append(results, "</varDec>")
	}
	tok3 := compileStatements(tok1)
	compileToken(tok3) // '}'
	results = append(results, "</subroutineBody>")

	results = append(results, "</subroutineDec>")
}

// parameterList = ((type varName) (',' type varName)* )?
func compileParameterList(now token) token {
	results = append(results, "<parameterList>")
	if now.str == RPAREN {
		results = append(results, "</parameterList>")
		return now
	}
	compileToken(now)         // type
	compileToken(nextToken()) // varName
	tok := nextToken()
	for tok.str == COMMA {
		compileToken(tok)         // ','
		compileToken(nextToken()) // type
		compileToken(nextToken()) // varName
		tok = nextToken()
	}
	results = append(results, "</parameterList>")
	return tok
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

func nextToken() token {
	position++
	return tokens[position]
}

func readNextToken() token {
	return tokens[position+1]
}
