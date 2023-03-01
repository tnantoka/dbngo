package parser

func Parse(yylex yyLexer) int {
	yylex.(*Lexer).Whitespace ^= 1 << '\n'

	return yyParse(yylex)
}
