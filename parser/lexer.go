package parser

import (
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner
	Statements []Statement
	Errors     []string
}

func (l *Lexer) Lex(lval *yySymType) int {
	token := int(l.Scan())
	literal := l.TokenText()
	switch token {
	case scanner.Int:
		token = NUMBER
	case scanner.Ident:
		switch literal {
		case "Paper":
			token = PAPER
		}
	case '\n':
		token = LF
	case scanner.EOF:
		token = 0
	}
	lval.token = Token{Token: token, Literal: l.TokenText()}

	return token
}

func (l *Lexer) Error(e string) {
	l.Errors = append(l.Errors, e)
}
