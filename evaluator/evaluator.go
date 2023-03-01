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

	parser.Parse(l)
	e.Errors = l.Errors

	if len(e.Errors) > 0 {
		return e.img
	}

	draw.Draw(e.img, e.img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{0, 0}, draw.Src)

	env := NewEnvironment()
	e.evalStatements(l.Statements, env)

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

func (e *Evaluator) evalStatements(statements []parser.Statement, env *Environment) {
	for _, statement := range statements {
		switch s := statement.(type) {
		case *parser.PaperStatement:
			e.evalPaperStatement(s, env)
		case *parser.PenStatement:
			e.evalPenStatement(s, env)
		case *parser.LineStatement:
			e.evalLineStatement(s, env)
		case *parser.SetStatement:
			e.evalSetStatement(s, env)
		case *parser.BlockStatement:
			e.evalStatements(s.Statements, env)
		case *parser.RepeatStatement:
			e.evalRepeatStatement(s, env)
		}
	}
}

func (e *Evaluator) evalPaperStatement(statement *parser.PaperStatement, env *Environment) {
	draw.Draw(e.img, e.img.Bounds(), &image.Uniform{e.evalColor(statement.Value, env)}, image.Point{0, 0}, draw.Src)
}

func (e *Evaluator) evalPenStatement(statement *parser.PenStatement, env *Environment) {
	e.color = e.evalColor(statement.Value, env)
}

func (e *Evaluator) evalLineStatement(statement *parser.LineStatement, env *Environment) {
	x1 := e.evalNumber(statement.X1, env)
	y1 := 100 - e.evalNumber(statement.Y1, env)
	x2 := e.evalNumber(statement.X2, env)
	y2 := 100 - e.evalNumber(statement.Y2, env)
	bresenham.DrawLine(e.img, x1, y1, x2, y2, e.color)
}

func (e *Evaluator) evalSetStatement(statement *parser.SetStatement, env *Environment) {
	switch s := statement.Value.(type) {
	case *parser.NumberExpression:
		env.Set(statement.Name, e.evalNumber(s, env))
	}
}

func (e *Evaluator) evalRepeatStatement(statement *parser.RepeatStatement, env *Environment) {
	for i := e.evalNumber(statement.From, env); i <= e.evalNumber(statement.To, env); i++ {
		env.Set(statement.Name, i)
		e.evalStatements(statement.Body.(*parser.BlockStatement).Statements, env)
	}
}

func (e *Evaluator) evalColor(expression parser.Expression, env *Environment) color.Color {
	switch exp := expression.(type) {
	case *parser.NumberExpression, *parser.IdentifierExpression:
		num := e.evalNumber(exp, env)
		col := uint8((100 - num) * 255 / 100)
		return color.RGBA{col, col, col, 255}
	}
	return color.RGBA{0, 0, 0, 0}
}

func (e *Evaluator) evalNumber(expression parser.Expression, env *Environment) int {
	switch exp := expression.(type) {
	case *parser.NumberExpression:
		num, _ := strconv.Atoi(exp.Literal)
		return num
	case *parser.IdentifierExpression:
		num, ok := env.Get(exp.Literal)
		if !ok {
			e.Errors = append(e.Errors, "Identifier not found: "+exp.Literal)
		}
		return num
	}
	return 0
}
