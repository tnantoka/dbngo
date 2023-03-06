package parser

import (
	"strings"
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
	case scanner.String:
		literal = strings.Trim(literal, "\"")
		token = STRING
	case scanner.Ident:
		switch literal {
		case "Paper":
			token = PAPER
		case "Pen":
			token = PEN
		case "Line":
			token = LINE
		case "Set":
			token = SET
		case "Repeat":
			token = REPEAT
		case "Same":
			if l.Peek() == '?' {
				token = SAME
				l.Next()
			}
		case "NotSame":
			if l.Peek() == '?' {
				token = NOTSAME
				l.Next()
			}
		case "Smaller":
			if l.Peek() == '?' {
				token = SMALLER
				l.Next()
			}
		case "NotSmaller":
			if l.Peek() == '?' {
				token = NOTSMALLER
				l.Next()
			}
		case "Command":
			token = COMMAND
		case "Load":
			token = LOAD
		default:
			token = IDENTIFIER
		}
	case '{':
		token = LBRACE
	case '}':
		token = RBRACE
	case '(':
		token = LPAREN
	case ')':
		token = RPAREN
	case '[':
		token = LBRACKET
	case ']':
		token = RBRACKET
	case '+', '-', '*', '/':
		token = OPERATOR
	case '\n':
		token = LF
	case scanner.EOF:
		token = 0
	}
	lval.token = Token{Token: token, Literal: literal}

	return token
}

func (l *Lexer) Error(e string) {
	l.Errors = append(l.Errors, e)
}
