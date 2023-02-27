package evaluator

import (
	"image/color"
	"strings"
	"testing"
)

func TestSyntax(t *testing.T) {
	tests := []struct {
		input    string
		expected [][]color.Color
	}{
		{
			"",
			[][]color.Color{
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			},
		},
		{
			"\n",
			[][]color.Color{
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			},
		},
		{
			"Paper 100\n",
			[][]color.Color{
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			},
		},
		{
			"Paper 100",
			[][]color.Color{
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			},
		},
		{
			"Paper 100\n\n",
			[][]color.Color{
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			},
		},
		{
			"Paper 100\nPaper 100\n",
			[][]color.Color{
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			},
		},
		{
			"Paper 100\n\nPaper 100\n\n\n",
			[][]color.Color{
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
				{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			},
		},
	}

	for i, test := range tests {
		e := New()
		e.length = 2
		pixels := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		for y := range pixels {
			for x := range pixels[y] {
				if pixels[y][x] != test.expected[y][x] {
					t.Errorf("test %d [%d][%d]: expected %v, got %v", i, y, x, test.expected, pixels)
				}
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
			"Paper 100 Paper 100\n",
			[]string{
				"syntax error",
			},
		},
	}

	for i, test := range tests {
		e := New()
		e.Eval(strings.NewReader(test.input))

		if len(e.Errors) != len(test.expected) {
			t.Errorf("test %d: expected %d errors, got %d", i, len(test.expected), len(e.Errors))
		}

		for j, err := range e.Errors {
			if err != test.expected[j] {
				t.Errorf("test %d: expected %s, got %s", i, test.expected[j], err)
			}
		}
	}
}

func TestPaper(t *testing.T) {
	tests := []struct {
		input    string
		expected [][]color.Color
	}{
		{
			"Paper 50",
			[][]color.Color{
				{color.RGBA{127, 127, 127, 255}, color.RGBA{127, 127, 127, 255}},
				{color.RGBA{127, 127, 127, 255}, color.RGBA{127, 127, 127, 255}},
			},
		},
		{
			"Paper 50\nPaper 10",
			[][]color.Color{
				{color.RGBA{25, 25, 25, 255}, color.RGBA{25, 25, 25, 255}},
				{color.RGBA{25, 25, 25, 255}, color.RGBA{25, 25, 25, 255}},
			},
		},
	}

	for i, test := range tests {
		e := New()
		e.length = 2
		pixels := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		for y := range pixels {
			for x := range pixels[y] {
				if pixels[y][x] != test.expected[y][x] {
					t.Errorf("test %d [%d][%d]: expected %v, got %v", i, y, x, test.expected, pixels)
				}
			}
		}
	}
}
