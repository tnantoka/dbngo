package evaluator

import (
	"image/color"
	"io"
	"strconv"

	"github.com/tnantoka/dbngo/parser"
)

const DEFAULT_LENGTH = 100

type Evaluator struct {
	length int
	Errors []string
	pixels [][]color.Color
	color  color.Color
}

func New() *Evaluator {
	return &Evaluator{length: DEFAULT_LENGTH, color: color.RGBA{0, 0, 0, 255}}
}

func (e *Evaluator) Eval(input io.Reader) [][]color.Color {
	l := new(parser.Lexer)
	l.Init(input)

	l.Whitespace ^= 1 << '\n'

	parser.Parse(l)
	e.Errors = l.Errors

	e.initPixels()
	e.evalStatements(l.Statements)

	return e.pixels
}

func (e *Evaluator) initPixels() {
	e.pixels = make([][]color.Color, e.length)
	for i := 0; i < e.length; i++ {
		e.pixels[i] = make([]color.Color, e.length)
	}

	for i := 0; i < e.length; i++ {
		for j := 0; j < e.length; j++ {
			e.pixels[i][j] = color.RGBA{255, 255, 255, 255}
		}
	}
}

func (e *Evaluator) evalStatements(statements []parser.Statement) {
	for _, statement := range statements {
		switch s := statement.(type) {
		case *parser.PaperStatement:
			e.evalPaperStatement(s)
		case *parser.PenStatement:
			e.evalPenStatement(s)
		case *parser.LineStatement:
			e.evalLineStatement(s)
		}
	}
}

func (e *Evaluator) evalPaperStatement(statement *parser.PaperStatement) {
	for i := 0; i < e.length; i++ {
		for j := 0; j < e.length; j++ {
			e.pixels[i][j] = evalColor(statement.Value)
		}
	}
}

func (e *Evaluator) evalPenStatement(statement *parser.PenStatement) {
	e.color = evalColor(statement.Value)
}

func (e *Evaluator) evalLineStatement(statement *parser.LineStatement) {
}

func evalColor(expression parser.Expression) color.Color {
	switch e := expression.(type) {
	case *parser.NumberExpression:
		num, _ := strconv.Atoi(e.Literal)
		col := uint8((100 - num) * 255 / 100)
		return color.RGBA{col, col, col, 255}
	}
	return color.RGBA{0, 0, 0, 0}
}
