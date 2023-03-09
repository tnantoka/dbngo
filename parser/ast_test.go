package parser

import (
	"testing"
	"text/scanner"
)

func TestPost(t *testing.T) {
	position := scanner.Position{
		Filename: "test",
		Line:     1,
		Column:   1,
	}
	token := Token{Position: position}

	expect := "test:1:1: "
	actual := token.Pos()
	if expect != actual {
		t.Errorf("expect %s, but got %s", expect, actual)
	}
}
