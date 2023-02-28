package evaluator

import (
	"image"
	"image/color"
	"io"
	"strconv"

	"github.com/StephaneBunel/bresenham"
	"github.com/tnantoka/dbngo/parser"
	"golang.org/x/image/draw"
)

const DEFAULT_LENGTH = 100

type Evaluator struct {
	length int
	Errors []string
	color  color.Color
	img    *image.RGBA
	Scale  int
}

func New() *Evaluator {
	return &Evaluator{length: DEFAULT_LENGTH, color: color.RGBA{0, 0, 0, 255}, Scale: 1}
}

func (e *Evaluator) Eval(input io.Reader) image.Image {
	e.img = image.NewRGBA(image.Rect(0, 0, e.length, e.length))

	l := new(parser.Lexer)
	l.Init(input)

	l.Whitespace ^= 1 << '\n'

	parser.Parse(l)
	e.Errors = l.Errors

	if len(e.Errors) > 0 {
		return e.img
	}

	draw.Draw(e.img, e.img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{0, 0}, draw.Src)

	e.evalStatements(l.Statements)

	return e.scale()
}

func (e *Evaluator) scale() image.Image {
	if e.Scale < 2 {
		return e.img
	}

	scaled := image.NewRGBA(image.Rect(0, 0, e.length*e.Scale, e.length*e.Scale))
	draw.CatmullRom.Scale(scaled, scaled.Bounds(), e.img, e.img.Bounds(), draw.Over, nil)

	return scaled
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
	draw.Draw(e.img, e.img.Bounds(), &image.Uniform{evalColor(statement.Value)}, image.Point{0, 0}, draw.Src)
}

func (e *Evaluator) evalPenStatement(statement *parser.PenStatement) {
	e.color = evalColor(statement.Value)
}

func (e *Evaluator) evalLineStatement(statement *parser.LineStatement) {
	x1 := evalNumber(statement.X1)
	y1 := 100 - evalNumber(statement.Y1)
	x2 := evalNumber(statement.X2)
	y2 := 100 - evalNumber(statement.Y2)
	bresenham.DrawLine(e.img, x1, y1, x2, y2, e.color)
}

func evalColor(expression parser.Expression) color.Color {
	switch e := expression.(type) {
	case *parser.NumberExpression:
		num := evalNumber(e)
		col := uint8((100 - num) * 255 / 100)
		return color.RGBA{col, col, col, 255}
	}
	return color.RGBA{0, 0, 0, 0}
}

func evalNumber(expression parser.Expression) int {
	switch e := expression.(type) {
	case *parser.NumberExpression:
		num, _ := strconv.Atoi(e.Literal)
		return num
	}
	return 0
}
