package parser

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected []Statement
	}{
		{
			input: "Paper 100\n",
			expected: []Statement{
				&PaperStatement{Value: &NumberExpression{Literal: "100"}},
			},
		},
		{
			input: "Pen (10 + 10)",
			expected: []Statement{
				&PenStatement{Value: &CalculateExpression{
					Left:     &NumberExpression{Literal: "10"},
					Operator: "+",
					Right:    &NumberExpression{Literal: "10"},
				}},
			},
		},
		{
			input: "Line 0 0 100 100",
			expected: []Statement{
				&LineStatement{
					X1: &NumberExpression{Literal: "0"},
					Y1: &NumberExpression{Literal: "0"},
					X2: &NumberExpression{Literal: "100"},
					Y2: &NumberExpression{Literal: "100"},
				},
			},
		},
		{
			input: "Set X 100",
			expected: []Statement{
				&SetStatement{
					Name:  "X",
					Value: &NumberExpression{Literal: "100"},
				},
			},
		},
		{
			input: "Set [1 2] 100",
			expected: []Statement{
				&DotStatement{
					X:     &NumberExpression{Literal: "1"},
					Y:     &NumberExpression{Literal: "2"},
					Value: &NumberExpression{Literal: "100"},
				},
			},
		},
		{
			input: "Repeat X 0 10 { Pen X }",
			expected: []Statement{
				&RepeatStatement{
					Name: "X",
					From: &NumberExpression{Literal: "0"},
					To:   &NumberExpression{Literal: "10"},
					Body: &BlockStatement{
						Statements: []Statement{
							&PenStatement{Value: &IdentifierExpression{Literal: "X"}},
						},
					},
				},
			},
		},
	}

	for i, test := range tests {
		l := new(Lexer)
		l.Init(strings.NewReader(test.input))

		Parse(l)

		if len(l.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, l.Errors)
		}

		if len(test.expected) != len(l.Statements) {
			t.Errorf("test %d: expected %d statements, got %d", i, len(test.expected), len(l.Statements))
		}

		for j, statement := range l.Statements {
			if statement.String() != test.expected[j].String() {
				t.Errorf("test %d: expected %s, got %s", i, test.expected[j], statement)
			}
		}
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: "Paper 100 Paper 100\n",
			expected: []string{
				"syntax error",
			},
		},
	}

	for i, test := range tests {
		l := new(Lexer)
		l.Init(strings.NewReader(test.input))

		Parse(l)

		if len(l.Errors) != len(test.expected) {
			t.Errorf("test %d: expected %d errors, got %d", i, len(test.expected), len(l.Errors))
		}

		for j, err := range l.Errors {
			if err != test.expected[j] {
				t.Errorf("test %d: expected %s, got %s", i, test.expected[j], err)
			}
		}
	}
}
