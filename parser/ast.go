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
