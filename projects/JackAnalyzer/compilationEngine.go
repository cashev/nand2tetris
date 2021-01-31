package main

var now = 0
var tokens []token

func compile(toks []token) string {
	tokens = toks

	return ""
}

func compileClass() []string {
	results := []string{"<class>"}
	str := compileToken(tokens[pos])
	results = append(results, str)
	pos++
	tok := tokens[pos]
	for tok.kind != "!" {

		pos++
		tok = tokens[pos]
	}

	results = append(results, "</class>")
	return results
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
