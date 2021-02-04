package main

var position int
var tokens []token

func compile(toks []token) []string {
	position = 0
	tokens = toks

	var results []string
	for position+1 < len(tokens) {
		tok := tokens[position]
		results = append(results, compileClass(tok)...)
	}
	return results
}

// class = 'class' className '{' classVarDec* subroutineDec* '}'
func compileClass(now token) []string {
	results := []string{"<class>"}
	results = append(results, compileToken(now))         // class
	results = append(results, compileToken(nextToken())) // className
	results = append(results, compileToken(nextToken())) // {
	tok := nextToken()
	for tok.str != RBRACE {
		if tok.str == "static" || tok.str == "field" {
			compileClassVarDec(tok)
		}
		if isSubroutine(tok.str) {
			subroutine := compileSubroutine(tok)
			results = append(results, subroutine...)
		}
		tok = nextToken()
	}
	results = append(results, compileToken(tok)) // }
	results = append(results, "</class>")
	return results
}

// classVarDec = ('static' | 'field') type varName (',' varName)* ';'
func compileClassVarDec(now token) []string {
	results := []string{"<classVarDec>"}
	results = append(results, compileToken(now))         // ('static' | 'field')
	results = append(results, compileToken(nextToken())) // type
	results = append(results, compileToken(nextToken())) // varName
	for tok := nextToken(); tok.str == COMMA; tok = nextToken() {
		results = append(results, compileToken(tok))         // ','
		results = append(results, compileToken(nextToken())) // varName
	}
	results = append(results, compileToken(nextToken())) // ';'
	results = append(results, "</classVarDec>")
	return results
}

func isSubroutine(str string) bool {
	return str == FUNCTION || str == METHOD || str == CONSTRUCTOR
}

// subroutineDec = ('construct' | 'function' | 'method')
//									('void' | type) subroutineName '(' parameterList ')'
//									subroutineBody
func compileSubroutine(now token) []string {
	results := []string{"<subroutineDec>"}
	results = append(results, compileToken(now))         // subroutineDec
	results = append(results, compileToken(nextToken())) // type
	results = append(results, compileToken(nextToken())) // subroutineName

	results = append(results, compileToken(nextToken())) // (
	paramList, tok := compileParameterList(nextToken())
	results = append(results, paramList...)
	results = append(results, compileToken(tok)) // )

	// subroutineBody = '{' varDec* statements '}'
	results = append(results, "<subroutineBody>")
	results = append(results, compileToken(nextToken())) // {
	tok1 := nextToken()
	for tok1.str == "var" {
		results = append(results, "<varDec>")
		// varDec = 'var' type varName (',' varName)* ';'
		results = append(results, compileToken(tok1))        // var
		results = append(results, compileToken(nextToken())) // type
		results = append(results, compileToken(nextToken())) // varName
		tok2 := nextToken()
		for tok2.str == COMMA {
			results = append(results, compileToken(tok2))        // ','
			results = append(results, compileToken(nextToken())) // varName
			tok2 = nextToken()
		}
		results = append(results, compileToken(tok2)) // ';'
		tok1 = nextToken()
		results = append(results, "</varDec>")
	}
	stmts, tok3 := compileStatements(tok1)
	results = append(results, stmts...)
	results = append(results, compileToken(tok3)) // }
	results = append(results, "</subroutineBody>")

	results = append(results, "</subroutineDec>")
	return results
}

// parameterList = ((type varName) (',' type varName)* )?
func compileParameterList(now token) ([]string, token) {
	results := []string{"<parameterList>"}
	if now.str == RPAREN {
		results = append(results, "</parameterList>")
		return results, now
	}
	results = append(results, compileToken(now))         // type
	results = append(results, compileToken(nextToken())) // varName
	tok := nextToken()
	for tok.str == COMMA {
		results = append(results, compileToken(tok))         // ','
		results = append(results, compileToken(nextToken())) // type
		results = append(results, compileToken(nextToken())) // varName
		tok = nextToken()
	}
	results = append(results, "</parameterList>")
	return results, tok
}

func compileVarDec() {

}

// statements = statement*
func compileStatements(now token) ([]string, token) {
	results := []string{"<statements>"}
	// statement = letStatement | ifStatement |
	// 						 whileStatement | doStatement |
	// 						 returnStatement
	loop := true
	tok := now
	for {
		switch tok.str {
		case LET:
			results = append(results, compileLet(tok)...)
		case IF:
			results = append(results, compileIf(tok)...)
		case WHILE:
			results = append(results, compileWhile(tok)...)
		case DO:
			results = append(results, compileDo(tok)...)
		case RETURN:
			results = append(results, compileReturn(tok)...)
		default:
			loop = false
		}
		if !loop {
			break
		}
		tok = nextToken()
	}
	results = append(results, "</statements>")
	return results, tok
}

// doStatement = 'do' subroutineCall ';'
func compileDo(now token) []string {
	results := []string{"<doStatement>"}
	results = append(results, compileToken(now))         // do
	results = append(results, compileToken(nextToken())) // identifier
	tok := nextToken()
	if tok.str == LPAREN {
		// subroutineCall = subroutineName '(' exprList ')'
		results = append(results, compileToken(tok)) // (
		exprList, next := compileExpressionList(nextToken())
		results = append(results, exprList...)
		results = append(results, compileToken(next)) // )
	}
	if tok.str == PERIOD {
		// subroutineCall = (className | varName) '.' subroutineName
		// 									'(' exprList ')'
		results = append(results, compileToken(tok))         // .
		results = append(results, compileToken(nextToken())) // identifier
		results = append(results, compileToken(nextToken())) // (
		exprList, next := compileExpressionList(nextToken())
		results = append(results, exprList...)
		results = append(results, compileToken(next)) // )
	}
	results = append(results, compileToken(nextToken())) // ;

	results = append(results, "</doStatement>")
	return results
}

// letstatement = 'let' varName ('[' expr ']')? '=' expr ';'
func compileLet(now token) []string {
	results := []string{"<letStatement>"}
	results = append(results, compileToken(now))         // let
	results = append(results, compileToken(nextToken())) // varName
	tok := nextToken()
	if tok.str == LBRACKET {
		results = append(results, compileToken(tok)) // [
		expr, tok1 := compileExpression(nextToken())
		results = append(results, expr...)
		results = append(results, compileToken(tok1)) // ]
		tok = nextToken()
	}
	results = append(results, compileToken(tok)) // =
	expr, tok := compileExpression(nextToken())
	results = append(results, expr...)
	results = append(results, compileToken(tok)) // ;

	results = append(results, "</letStatement>")
	return results
}

// whileStatement = 'while' '(' expr ')' '{' statements '}'
func compileWhile(now token) []string {
	results := []string{"<whileStatement>"}
	results = append(results, compileToken(now))         // while
	results = append(results, compileToken(nextToken())) // (
	expr, tok := compileExpression(nextToken())
	results = append(results, expr...)
	results = append(results, compileToken(tok))         // )
	results = append(results, compileToken(nextToken())) // {
	stmts, tok := compileStatements(nextToken())
	results = append(results, stmts...)
	results = append(results, compileToken(tok)) // }
	results = append(results, "</whileStatement>")
	return results
}

// returnStatement = 'return' expr? ';'
func compileReturn(now token) []string {
	results := []string{"<returnStatement>"}
	results = append(results, compileToken(now)) // return
	tok := nextToken()
	if tok.str != SEMICOLON {
		expr, tok1 := compileExpression(tok)
		results = append(results, expr...)
		tok = tok1
	}
	results = append(results, compileToken(tok)) // ';'
	results = append(results, "</returnStatement>")
	return results
}

// ifStatement = 'if' '(' expr ')' '{' statements '}'
// 							 ('else' '{' statements '}')?
func compileIf(now token) []string {
	results := []string{"<ifStatement>"}
	results = append(results, compileToken(now)) // if
	results = append(results, "</ifStatement>")
	return results
}

// expr = term (op term)*
func compileExpression(now token) ([]string, token) {
	results := []string{"<expression>"}
	results = append(results, compileTerm(now)...) // term
	tok := nextToken()
	for isOp(tok.str) {
		results = append(results, compileToken(tok))           // op
		results = append(results, compileTerm(nextToken())...) // term
		tok = nextToken()
	}
	results = append(results, "</expression>")
	return results, tok
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
func compileTerm(now token) []string {
	results := []string{"<term>"}
	if now.kind == IDENTIFIER {
		results = append(results, compileToken(now))
		if readNextToken().str == LBRACKET {
			// varName '[' expr ']'
			results = append(results, compileToken(nextToken())) // [
			expr, tok := compileExpression(nextToken())
			results = append(results, expr...)
			results = append(results, compileToken(tok)) // ]
		}
		if readNextToken().str == LPAREN {
			// subroutineCall = subroutineName '(' exprList ')'
			results = append(results, compileToken(nextToken())) // (
			exprList, next := compileExpressionList(nextToken())
			results = append(results, exprList...)
			results = append(results, compileToken(next)) // )

		}
		if readNextToken().str == PERIOD {
			// subroutineCall = (className | varName) '.' subroutineName
			// 									'(' exprList ')'
			results = append(results, compileToken(nextToken())) // .
			results = append(results, compileToken(nextToken())) // subroutineName
			results = append(results, compileToken(nextToken())) // (
			exprList, next := compileExpressionList(nextToken())
			results = append(results, exprList...)
			results = append(results, compileToken(next)) // )
		}

	} else if isUnaryOp(now.str) {
		results = append(results, compileToken(now))           // uparyOp
		results = append(results, compileTerm(nextToken())...) // term
	} else if now.str == LPAREN {
		// '(' expr ')'
		results = append(results, compileToken(now)) // (
		expr, tok := compileExpression(nextToken())
		results = append(results, expr...)
		results = append(results, compileToken(tok)) // (
	} else {
		results = append(results, compileToken(now))
	}
	results = append(results, "</term>")
	return results
}

// exprList = (expr (',' expr)* )?
func compileExpressionList(now token) ([]string, token) {
	results := []string{"<expressionList>"}
	if now.str == RPAREN {
		results = append(results, "</expressionList>")
		return results, now
	}
	expr, tok := compileExpression(now)
	results = append(results, expr...)

	for tok.str == COMMA {
		results = append(results, compileToken(tok)) // ,
		expr, tok = compileExpression(now)
		results = append(results, expr...)
	}
	results = append(results, "</expressionList>")
	return results, tok
}

func compileToken(tok token) string {
	kind := string(tok.kind)
	str := tok.str
	return "<" + kind + "> " + str + " </" + kind + ">"
}

func compileTokenList(toks []token) []string {
	results := []string{"<tokens>"}
	for _, tok := range toks {
		str := compileToken(tok)
		results = append(results, str)
	}
	results = append(results, "</tokens>")
	return results
}

func nextToken() token {
	position++
	return tokens[position]
}

func readNextToken() token {
	return tokens[position+1]
}
