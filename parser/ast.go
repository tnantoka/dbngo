package parser

import (
	"strings"
)

type Token struct {
	Token   int
	Literal string
}

type Expression interface {
	String() string
}

type IntegerExpression struct {
	Literal string
}

func (ne *IntegerExpression) String() string {
	return ne.Literal
}

type IdentifierExpression struct {
	Literal string
}

func (ie *IdentifierExpression) String() string {
	return ie.Literal
}

type CalculateExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ce *CalculateExpression) String() string {
	return ce.Left.String() + " " + ce.Operator + " " + ce.Right.String()
}

type CallNumberExpression struct {
	Name      string
	Arguments []Expression
}

func (ce *CallNumberExpression) String() string {
	var out string = "<"
	out += ce.Name
	for _, a := range ce.Arguments {
		out += " " + a.String()
	}
	out += ">"
	return out
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

type SetStatement struct {
	Name  string
	Value Expression
}

func (ss *SetStatement) String() string {
	return "Set " + ss.Name + " " + ss.Value.String()
}

type DotStatement struct {
	X     Expression
	Y     Expression
	Value Expression
}

func (ds *DotStatement) String() string {
	return "Set [" + ds.X.String() + " " + ds.Y.String() + "] " + ds.Value.String()
}

type CopyStatement struct {
	Name string
	X    Expression
	Y    Expression
}

func (cs *CopyStatement) String() string {
	return "Set " + cs.Name + " [" + cs.X.String() + " " + cs.Y.String() + "]"
}

type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) String() string {
	var out string
	out += "{\n"
	for _, s := range bs.Statements {
		out += s.String() + "\n"
	}
	out += "}"
	return out
}

type RepeatStatement struct {
	Name string
	From Expression
	To   Expression
	Body Statement
}

func (rs *RepeatStatement) String() string {
	return "Repeat " + rs.Body.String()
}

type SameStatement struct {
	Left  Expression
	Right Expression
	Body  Statement
}

func (ss *SameStatement) String() string {
	return "Same? " + ss.Left.String() + " " + ss.Right.String() + " " + ss.Body.String()
}

type NotSameStatement struct {
	Left  Expression
	Right Expression
	Body  Statement
}

func (ns *NotSameStatement) String() string {
	return "NotSame? " + ns.Left.String() + " " + ns.Right.String() + " " + ns.Body.String()
}

type SmallerStatement struct {
	Left  Expression
	Right Expression
	Body  Statement
}

func (ss *SmallerStatement) String() string {
	return "Smaller? " + ss.Left.String() + " " + ss.Right.String() + " " + ss.Body.String()
}

type NotSmallerStatement struct {
	Left  Expression
	Right Expression
	Body  Statement
}

func (ns *NotSmallerStatement) String() string {
	return "NotSmaller? " + ns.Left.String() + " " + ns.Right.String() + " " + ns.Body.String()
}

type DefineCommandStatement struct {
	Name       string
	Body       Statement
	Parameters []string
}

func (fs *DefineCommandStatement) String() string {
	return "Command " + fs.Name + " " + strings.Join(fs.Parameters, " ") + " " + fs.Body.String()
}

type CallCommandStatement struct {
	Name      string
	Arguments []Expression
}

func (cs *CallCommandStatement) String() string {
	var out string
	out += cs.Name
	for _, a := range cs.Arguments {
		out += " " + a.String()
	}
	return out
}

type LoadStatement struct {
	Path string
}

func (ls *LoadStatement) String() string {
	return "Load \"" + ls.Path + "\""
}

type DefineNumberStatement struct {
	Name       string
	Body       Statement
	Parameters []string
}

func (fs *DefineNumberStatement) String() string {
	return "Number " + fs.Name + " " + strings.Join(fs.Parameters, " ") + " " + fs.Body.String()
}

type ValueStatement struct {
	Result Expression
}

func (vs *ValueStatement) String() string {
	return "Value " + vs.Result.String()
}
