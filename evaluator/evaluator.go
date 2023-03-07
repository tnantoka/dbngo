package evaluator

import (
	"image"
	"image/color"
	"io"
	"os"
	"strconv"

	"github.com/StephaneBunel/bresenham"
	"github.com/tnantoka/dbngo/parser"
	"golang.org/x/image/draw"
)

const DEFAULT_LENGTH = 100

type Evaluator struct {
	length    int
	Errors    []string
	color     color.Color
	img       *image.RGBA
	Scale     int
	Directory string
}

func New() *Evaluator {
	return &Evaluator{length: DEFAULT_LENGTH, color: color.RGBA{0, 0, 0, 255}, Scale: 1, Directory: ""}
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
		e.evalStatement(statement, env)
	}
}

func (e *Evaluator) evalStatement(statement parser.Statement, env *Environment) {
	switch s := statement.(type) {
	case *parser.PaperStatement:
		e.evalPaperStatement(s, env)
	case *parser.PenStatement:
		e.evalPenStatement(s, env)
	case *parser.LineStatement:
		e.evalLineStatement(s, env)
	case *parser.SetStatement:
		e.evalSetStatement(s, env)
	case *parser.DotStatement:
		e.evalDotStatement(s, env)
	case *parser.CopyStatement:
		e.evalCopyStatement(s, env)
	case *parser.BlockStatement:
		e.evalStatements(s.Statements, env)
	case *parser.RepeatStatement:
		e.evalRepeatStatement(s, env)
	case *parser.SameStatement:
		e.evalSameStatement(s, env)
	case *parser.NotSameStatement:
		e.evalNotSameStatement(s, env)
	case *parser.SmallerStatement:
		e.evalSmallerStatement(s, env)
	case *parser.NotSmallerStatement:
		e.evalNotSmallerStatement(s, env)
	case *parser.DefineCommandStatement:
		e.evalDefineCommandStatement(s, env)
	case *parser.CallCommandStatement:
		e.evalCallCommandStatement(s, env)
	case *parser.LoadStatement:
		e.evalLoadStatement(s, env)
	case *parser.DefineNumberStatement:
		e.evalDefineNumberStatement(s, env)
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
	env.Set(statement.Name, e.evalNumber(statement.Value, env))
}

func (e *Evaluator) evalDotStatement(statement *parser.DotStatement, env *Environment) {
	x := e.evalNumber(statement.X, env)
	y := 100 - e.evalNumber(statement.Y, env)
	e.img.Set(x, y, e.color)
}

func (e *Evaluator) evalCopyStatement(statement *parser.CopyStatement, env *Environment) {
	name := statement.Name
	x := e.evalNumber(statement.X, env)
	y := 100 - e.evalNumber(statement.Y, env)
	r, _, _, _ := e.img.At(x, y).RGBA()
	env.Set(name, int(100-r*100/65535))
}

func (e *Evaluator) evalRepeatStatement(statement *parser.RepeatStatement, env *Environment) {
	for i := e.evalNumber(statement.From, env); i <= e.evalNumber(statement.To, env); i++ {
		env.Set(statement.Name, i)
		e.evalStatements(statement.Body.(*parser.BlockStatement).Statements, env)
	}
}

func (e *Evaluator) evalSameStatement(statement *parser.SameStatement, env *Environment) {
	left := e.evalNumber(statement.Left, env)
	right := e.evalNumber(statement.Right, env)
	if left == right {
		e.evalStatements(statement.Body.(*parser.BlockStatement).Statements, env)
	}
}

func (e *Evaluator) evalNotSameStatement(statement *parser.NotSameStatement, env *Environment) {
	left := e.evalNumber(statement.Left, env)
	right := e.evalNumber(statement.Right, env)
	if left != right {
		e.evalStatements(statement.Body.(*parser.BlockStatement).Statements, env)
	}
}

func (e *Evaluator) evalSmallerStatement(statement *parser.SmallerStatement, env *Environment) {
	left := e.evalNumber(statement.Left, env)
	right := e.evalNumber(statement.Right, env)
	if left < right {
		e.evalStatements(statement.Body.(*parser.BlockStatement).Statements, env)
	}
}

func (e *Evaluator) evalNotSmallerStatement(statement *parser.NotSmallerStatement, env *Environment) {
	left := e.evalNumber(statement.Left, env)
	right := e.evalNumber(statement.Right, env)
	if left >= right {
		e.evalStatements(statement.Body.(*parser.BlockStatement).Statements, env)
	}
}

func (e *Evaluator) evalDefineCommandStatement(statement *parser.DefineCommandStatement, env *Environment) {
	env.Set(statement.Name, statement)
}

func (e *Evaluator) evalCallCommandStatement(statement *parser.CallCommandStatement, env *Environment) {
	fun, ok := env.Get(statement.Name)
	if !ok {
		e.Errors = append(e.Errors, "Command not found: "+statement.Name)
		return
	}
	funStatement := fun.(*parser.DefineCommandStatement)
	newEnv := NewEnclosedEnvironment(env)
	for i, arg := range statement.Arguments {
		newEnv.Set(funStatement.Parameters[i], e.evalNumber(arg, env))
	}
	e.evalStatements(funStatement.Body.(*parser.BlockStatement).Statements, newEnv)
}

func (e *Evaluator) evalLoadStatement(statement *parser.LoadStatement, env *Environment) {
	file, err := os.Open(e.Directory + "/" + statement.Path)
	if err != nil {
		e.Errors = append(e.Errors, err.Error())
		return
	}

	l := new(parser.Lexer)
	l.Init(file)

	parser.Parse(l)
	e.Errors = append(e.Errors, l.Errors...)

	if len(l.Errors) > 0 {
		return
	}

	e.evalStatements(l.Statements, env)
}

func (e *Evaluator) evalDefineNumberStatement(statement *parser.DefineNumberStatement, env *Environment) {
	env.Set(statement.Name, statement)
}

func (e *Evaluator) evalColor(expression parser.Expression, env *Environment) color.Color {
	switch exp := expression.(type) {
	case *parser.IntegerExpression, *parser.IdentifierExpression, *parser.CalculateExpression, *parser.CallNumberExpression:
		num := e.evalNumber(exp, env)
		col := uint8((100 - num) * 255 / 100)
		return color.RGBA{col, col, col, 255}
	}
	return color.RGBA{0, 0, 0, 0}
}

func (e *Evaluator) evalNumber(expression parser.Expression, env *Environment) int {
	switch exp := expression.(type) {
	case *parser.IntegerExpression:
		num, _ := strconv.Atoi(exp.Literal)
		return num
	case *parser.IdentifierExpression:
		num, ok := env.Get(exp.Literal)
		if !ok || num == nil {
			e.Errors = append(e.Errors, "Identifier not found: "+exp.Literal)
			return 0
		}
		return num.(int)
	case *parser.CalculateExpression:
		left := e.evalNumber(exp.Left, env)
		right := e.evalNumber(exp.Right, env)
		switch exp.Operator {
		case "+":
			return left + right
		case "-":
			return left - right
		case "*":
			return left * right
		case "/":
			return left / right
		}
	case *parser.CallNumberExpression:
		fun, ok := env.Get(exp.Name)
		if !ok {
			e.Errors = append(e.Errors, "Number not found: "+exp.Name)
			return 0
		}
		funStatement := fun.(*parser.DefineNumberStatement)
		newEnv := NewEnclosedEnvironment(env)
		for i, arg := range exp.Arguments {
			newEnv.Set(funStatement.Parameters[i], e.evalNumber(arg, env))
		}
		for _, s := range funStatement.Body.(*parser.BlockStatement).Statements {
			vs, ok := s.(*parser.ValueStatement)
			if ok {
				return e.evalNumber(vs.Result, newEnv)
			} else {
				e.evalStatement(s, newEnv)
			}
		}
	}
	return 0
}
