package evaluator

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/tnantoka/dbngo/parser"
)

func TestSyntax(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"",
			"white.png",
		},
		{
			"\n",
			"white.png",
		},
		{
			"Paper 100",
			"black.png",
		},
		{
			"Paper 100\n",
			"black.png",
		},
		{
			"Paper 100\n\n",
			"black.png",
		},
		{
			"Paper 100\nPaper 100\n",
			"black.png",
		},
		{
			"Paper 100\n\nPaper 100\n\n\n",
			"black.png",
		},
	}

	for i, test := range tests {
		e := New()
		img := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		actual := imageToBytes(t, img)
		expected := readBytes(t, "../testdata/"+test.expected)

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
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
		{
			"Paper X\n",
			[]string{
				"Identifier not found: X",
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

func TestScale(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"",
			"scaled.png",
		},
	}

	for i, test := range tests {
		e := New()
		e.Scale = 2
		img := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		actual := imageToBytes(t, img)
		expected := readBytes(t, "../testdata/"+test.expected)

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
		}
	}
}

func TestPaper(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"Paper 50",
			"gray.png",
		},
		{
			"Paper 50\nPaper 10",
			"lightgray.png",
		},
	}

	for i, test := range tests {
		e := New()
		img := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		actual := imageToBytes(t, img)
		expected := readBytes(t, "../testdata/"+test.expected)

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
		}
	}
}

func TestPen(t *testing.T) {
	tests := []struct {
		input    string
		expected color.Color
	}{
		{
			"Pen 50",
			color.RGBA{127, 127, 127, 255},
		},
		{
			"Pen 50\nPen 10",
			color.RGBA{229, 229, 229, 255},
		},
	}

	for i, test := range tests {
		e := New()
		e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		if e.color != test.expected {
			t.Errorf("test %d: expected %v, got %v", i, test.expected, e.color)
		}
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		input    string
		expected color.Color
	}{
		{
			"Set X 50\nPen X",
			color.RGBA{127, 127, 127, 255},
		},
	}

	for i, test := range tests {
		e := New()
		e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		if e.color != test.expected {
			t.Errorf("test %d: expected %v, got %v", i, test.expected, e.color)
		}
	}
}

func TestDot(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"Set [1 1] 100",
			"dot.png",
		},
	}

	for i, test := range tests {
		e := New()
		img := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		actual := imageToBytes(t, img)
		expected := readBytes(t, "../testdata/"+test.expected)

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
		}
	}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		input    string
		expected color.Color
	}{
		{
			"Set X [1 1]\nPen X",
			color.RGBA{255, 255, 255, 255},
		},
	}

	for i, test := range tests {
		e := New()
		e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		if e.color != test.expected {
			t.Errorf("test %d: expected %v, got %v", i, test.expected, e.color)
		}
	}
}

func TestEvalColor(t *testing.T) {
	tests := []struct {
		input    fmt.Stringer
		expected color.Color
	}{
		{
			nil,
			color.RGBA{0, 0, 0, 0},
		},
		{
			&parser.NumberExpression{Literal: "10"},
			color.RGBA{229, 229, 229, 255},
		},
	}

	for i, test := range tests {
		e := New()
		evaluated := e.evalColor(test.input, NewEnvironment())
		if evaluated != test.expected {
			t.Errorf("test %d: expected %v, got %v", i, test.expected, evaluated)
		}
	}
}

func TestEvalNumber(t *testing.T) {
	tests := []struct {
		input    fmt.Stringer
		expected int
	}{
		{
			nil,
			0,
		},
		{
			&parser.NumberExpression{Literal: "10"},
			10,
		},
	}

	for i, test := range tests {
		e := New()
		evaluated := e.evalNumber(test.input, NewEnvironment())
		if evaluated != test.expected {
			t.Errorf("test %d: expected %v, got %v", i, test.expected, evaluated)
		}
	}
}

func TestLine(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"Line 0 0 100 100",
			"diagonal.png",
		},
	}

	for i, test := range tests {
		e := New()
		img := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		actual := imageToBytes(t, img)
		expected := readBytes(t, "../testdata/"+test.expected)

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
		}
	}
}

func TestBlock(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"{ Line 0 0 100 100 }",
			"diagonal.png",
		},
	}

	for i, test := range tests {
		e := New()
		img := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		actual := imageToBytes(t, img)
		expected := readBytes(t, "../testdata/"+test.expected)

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
		}
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"Repeat A 10 20 { Line A 10 A 20 }",
			"square.png",
		},
	}

	for i, test := range tests {
		e := New()
		img := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		actual := imageToBytes(t, img)
		expected := readBytes(t, "../testdata/"+test.expected)

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
		}
	}
}

func TestSame(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"Same? 1 1 { Line 0 0 100 100 }",
			"diagonal.png",
		},
	}

	for i, test := range tests {
		e := New()
		img := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		actual := imageToBytes(t, img)
		expected := readBytes(t, "../testdata/"+test.expected)

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
		}
	}
}

func TestNotSame(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"NotSame? 0 1 { Line 0 0 100 100 }",
			"diagonal.png",
		},
	}

	for i, test := range tests {
		e := New()
		img := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		actual := imageToBytes(t, img)
		expected := readBytes(t, "../testdata/"+test.expected)

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
		}
	}
}
func TestCalculate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"Paper (10 + 40)",
			"gray.png",
		},
		{
			"Paper (60 - 10)",
			"gray.png",
		},
		{
			"Paper (25 * 2)",
			"gray.png",
		},
		{
			"Paper (100 / 2)",
			"gray.png",
		},
	}

	for i, test := range tests {
		e := New()
		img := e.Eval(strings.NewReader(test.input))

		if len(e.Errors) > 0 {
			t.Errorf("test %d: expected no errors, got %v", i, e.Errors)
		}

		actual := imageToBytes(t, img)
		expected := readBytes(t, "../testdata/"+test.expected)

		if !bytes.Equal(actual, expected) {
			t.Errorf("test %d: expected %v, but got %v", i, expected, actual)
		}
	}
}

func imageToBytes(t *testing.T, img image.Image) []byte {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		t.Fatalf("failed to encode image: %s", err)
	}

	debugFile, err := os.Create("../tmp/debug.png")
	if err != nil {
		t.Fatalf("failed to create debug file: %s", err)
	}
	defer debugFile.Close()

	if err := png.Encode(debugFile, img); err != nil {
		t.Fatalf("failed to encode debug image: %s", err)
	}

	return buf.Bytes()
}

func readBytes(t *testing.T, path string) []byte {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("failed to read file: %s", err)
	}

	return bytes
}
