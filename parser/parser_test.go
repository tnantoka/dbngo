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
				&PaperStatement{Value: &IntegerExpression{Literal: "100"}},
			},
		},
		{
			input: "Pen (10 + 10)",
			expected: []Statement{
				&PenStatement{Value: &CalculateExpression{
					Left:     &IntegerExpression{Literal: "10"},
					Operator: "+",
					Right:    &IntegerExpression{Literal: "10"},
				}},
			},
		},
		{
			input: "Line 0 0 100 100",
			expected: []Statement{
				&LineStatement{
					X1: &IntegerExpression{Literal: "0"},
					Y1: &IntegerExpression{Literal: "0"},
					X2: &IntegerExpression{Literal: "100"},
					Y2: &IntegerExpression{Literal: "100"},
				},
			},
		},
		{
			input: "Set X 100",
			expected: []Statement{
				&SetStatement{
					Name:  "X",
					Value: &IntegerExpression{Literal: "100"},
				},
			},
		},
		{
			input: "Set [1 2] 100",
			expected: []Statement{
				&DotStatement{
					X:     &IntegerExpression{Literal: "1"},
					Y:     &IntegerExpression{Literal: "2"},
					Value: &IntegerExpression{Literal: "100"},
				},
			},
		},
		{
			input: "Set X [1 2]",
			expected: []Statement{
				&CopyStatement{
					Name: "X",
					X:    &IntegerExpression{Literal: "1"},
					Y:    &IntegerExpression{Literal: "2"},
				},
			},
		},
		{
			input: "Repeat X 0 10 { Pen X }",
			expected: []Statement{
				&RepeatStatement{
					Name: "X",
					From: &IntegerExpression{Literal: "0"},
					To:   &IntegerExpression{Literal: "10"},
					Body: &BlockStatement{
						Statements: []Statement{
							&PenStatement{Value: &IdentifierExpression{Token: Token{Literal: "X"}}},
						},
					},
				},
			},
		},
		{
			input: "Same? 0 10 { Pen X }",
			expected: []Statement{
				&SameStatement{
					Left:  &IntegerExpression{Literal: "0"},
					Right: &IntegerExpression{Literal: "10"},
					Body: &BlockStatement{
						Statements: []Statement{
							&PenStatement{Value: &IdentifierExpression{Token: Token{Literal: "X"}}},
						},
					},
				},
			},
		},
		{
			input: "NotSame? 0 10 { Pen X }",
			expected: []Statement{
				&NotSameStatement{
					Left:  &IntegerExpression{Literal: "0"},
					Right: &IntegerExpression{Literal: "10"},
					Body: &BlockStatement{
						Statements: []Statement{
							&PenStatement{Value: &IdentifierExpression{Token: Token{Literal: "X"}}},
						},
					},
				},
			},
		},
		{
			input: "Smaller? 0 10 { Pen X }",
			expected: []Statement{
				&SmallerStatement{
					Left:  &IntegerExpression{Literal: "0"},
					Right: &IntegerExpression{Literal: "10"},
					Body: &BlockStatement{
						Statements: []Statement{
							&PenStatement{Value: &IdentifierExpression{Token: Token{Literal: "X"}}},
						},
					},
				},
			},
		},
		{
			input: "NotSmaller? 0 10 { Pen X }",
			expected: []Statement{
				&NotSmallerStatement{
					Left:  &IntegerExpression{Literal: "0"},
					Right: &IntegerExpression{Literal: "10"},
					Body: &BlockStatement{
						Statements: []Statement{
							&PenStatement{Value: &IdentifierExpression{Token: Token{Literal: "X"}}},
						},
					},
				},
			},
		},
		{
			input: "Command Test X { Pen X }",
			expected: []Statement{
				&DefineCommandStatement{
					Name:       "Test",
					Parameters: []string{"X"},
					Body: &BlockStatement{
						Statements: []Statement{
							&PenStatement{Value: &IdentifierExpression{Token: Token{Literal: "X"}}},
						},
					},
				},
			},
		},
		{
			input: "Test 1",
			expected: []Statement{
				&CallCommandStatement{
					Token: Token{Literal: "Test"},
					Arguments: []Expression{
						&IntegerExpression{Literal: "1"},
					},
				},
			},
		},
		{
			input: "Load \"a.dbn\"",
			expected: []Statement{
				&LoadStatement{
					Token{Literal: "a.dbn"},
				},
			},
		},
		{
			input: "Number Test A { Value 1 }\nPaper <Test 1>",
			expected: []Statement{
				&DefineNumberStatement{
					Name:       "Test",
					Parameters: []string{"A"},
					Body: &BlockStatement{
						Statements: []Statement{
							&ValueStatement{
								Result: &IntegerExpression{Literal: "1"},
							},
						},
					},
				},
				&PaperStatement{
					Value: &CallNumberExpression{
						Token: Token{Literal: "Test"},
						Arguments: []Expression{
							&IntegerExpression{Literal: "1"},
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
				"test.dbn:1:11: syntax error",
			},
		},
	}

	for i, test := range tests {
		l := new(Lexer)
		l.Filename = "test.dbn"
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
