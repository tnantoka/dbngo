package parser

type Token struct {
	Token   int
	Literal string
}

type Expression interface {
	String() string
}

type NumberExpression struct {
	Literal string
}

func (ne *NumberExpression) String() string {
	return ne.Literal
}

type Statement interface {
	String() string
}

type PaperStatement struct {
	Value Expression
}

func (ps *PaperStatement) String() string {
	return "Paper " + ps.Value.String()
}

type PenStatement struct {
	Value Expression
}

func (ps *PenStatement) String() string {
	return "Pen " + ps.Value.String()
}

type LineStatement struct {
	X1 Expression
	Y1 Expression
	X2 Expression
	Y2 Expression
}

func (ls *LineStatement) String() string {
	return "Line " + ls.X1.String() + " " + ls.Y1.String() + " " + ls.X2.String() + " " + ls.Y2.String()
}
