package parser

import (
	"fmt"
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
		token = INTEGER
	case scanner.String:
		literal = strings.Trim(literal, "\"")
		token = STRING
	case scanner.Ident:
		switch literal {
		case "Paper", "paper":
			token = PAPER
		case "Pen", "pen":
			token = PEN
		case "Line", "line":
			token = LINE
		case "Set", "set":
			token = SET
		case "Repeat", "repeat":
			token = REPEAT
		case "Same", "same":
			if l.Peek() == '?' {
				token = SAME
				l.Next()
			}
		case "NotSame", "notsame":
			if l.Peek() == '?' {
				token = NOTSAME
				l.Next()
			}
		case "Smaller", "smaller":
			if l.Peek() == '?' {
				token = SMALLER
				l.Next()
			}
		case "NotSmaller", "notsmaller":
			if l.Peek() == '?' {
				token = NOTSMALLER
				l.Next()
			}
		case "Command", "command":
			token = COMMAND
		case "Load", "load":
			token = LOAD
		case "Number", "number":
			token = NUMBER
		case "Value", "value":
			token = VALUE
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
	case '<':
		token = LT
	case '>':
		token = GT
	case '\n':
		token = LF
	case scanner.EOF:
		token = 0
	}
	lval.token = Token{Token: token, Literal: literal, Position: l.Pos()}

	return token
}

func (l *Lexer) Error(e string) {
	l.Errors = append(l.Errors, fmt.Sprintf("%s:%d:%d: %s", l.Filename, l.Line, l.Column, e))
}
