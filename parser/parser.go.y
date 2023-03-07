%{
package parser
%}

%union{
    statements []Statement
    statement Statement
    expression Expression
    token Token
    parameters []string
    arguments []Expression
}

%type<statements> statements body

%type<statement> statement command
%type<statement> paper pen line set dot copy repeat same notsame smaller notsmaller definecommand callcommand load definenumber value
%type<statement> block

%type<expression> expression

%type<parameters> parameters
%type<arguments> arguments

%token<token> INTEGER LF IDENTIFIER OPERATOR
%token<token> PAPER PEN LINE SET REPEAT SAME NOTSAME SMALLER NOTSMALLER COMMAND LOAD NUMBER VALUE
%token<token> LBRACE RBRACE LPAREN RPAREN LBRACKET RBRACKET LT GT
%token<token> STRING

%%

statements
    : /* empty file */
    {
        $$ = nil
        yylex.(*Lexer).Statements = $$
    }
    | command /* no newline at end of file */
    {
        $$ = []Statement{$1}
        yylex.(*Lexer).Statements = $$
    }
    | statement statements
    {
        $$ = append([]Statement{$1}, $2...)
        yylex.(*Lexer).Statements = $$
    }

block
    : LBRACE body RBRACE
    {
        $$ = &BlockStatement{Statements: $2}
    }

body
    : /* empty block */
    {
        $$ = []Statement{}
    }
    | command /* no newline at end of block */
    {
        $$ = []Statement{$1}
    }
    | statement body
    {
        $$ = append([]Statement{$1}, $2...)
    }

statement
    : command LF
    | LF { $$ = nil }

command
    : paper
    | pen
    | line
    | set
    | dot
    | copy
    | block
    | repeat
    | same
    | notsame
    | smaller
    | notsmaller
    | definecommand
    | callcommand
    | load
    | definenumber
    | value

paper
    : PAPER expression
    {
        $$ = &PaperStatement{Value: $2}
    }

pen
    : PEN expression
    {
        $$ = &PenStatement{Value: $2}
    }

line
    : LINE expression expression expression expression
    {
        $$ = &LineStatement{X1: $2, Y1: $3, X2: $4, Y2: $5}
    }

set
    : SET IDENTIFIER expression
    {
        $$ = &SetStatement{Name: $2.Literal, Value: $3}
    }

dot
    : SET LBRACKET expression expression RBRACKET expression
    {
        $$ = &DotStatement{X: $3, Y: $4, Value: $6}
    }

copy
    : SET IDENTIFIER LBRACKET expression expression RBRACKET
    {
        $$ = &CopyStatement{Name: $2.Literal, X: $4, Y: $5}
    }

repeat
    : REPEAT IDENTIFIER expression expression block
    {
        $$ = &RepeatStatement{Name: $2.Literal, From: $3, To: $4, Body: $5}
    }

same
    : SAME expression expression block
    {
        $$ = &SameStatement{Left: $2, Right: $3, Body: $4}
    }

notsame
    : NOTSAME expression expression block
    {
        $$ = &NotSameStatement{Left: $2, Right: $3, Body: $4}
    }

smaller
    : SMALLER expression expression block
    {
        $$ = &SmallerStatement{Left: $2, Right: $3, Body: $4}
    }

notsmaller
    : NOTSMALLER expression expression block
    {
        $$ = &NotSmallerStatement{Left: $2, Right: $3, Body: $4}
    }

definecommand
    : COMMAND IDENTIFIER parameters block
    {
        $$ = &DefineCommandStatement{Name: $2.Literal, Parameters: $3, Body: $4}
    }

parameters
    : /* empty arguments */
    {
        $$ = []string{}
    }
    | IDENTIFIER parameters
    {
        $$ = append([]string{$1.Literal}, $2...)
    }

callcommand
    : IDENTIFIER arguments
    {
        $$ = &CallCommandStatement{Name: $1.Literal, Arguments: $2}
    }

definenumber
    : NUMBER IDENTIFIER parameters block
    {
        $$ = &DefineNumberStatement{Name: $2.Literal, Parameters: $3, Body: $4}
    }

arguments
    : /* empty arguments */
    {
        $$ = []Expression{}
    }
    | expression arguments
    {
        $$ = append([]Expression{$1}, $2...)
    }

load
    : LOAD STRING
    {
        $$ = &LoadStatement{Path: $2.Literal}
    }

value
   : VALUE expression
   {
       $$ = &ValueStatement{Result: $2}
   } 

expression
    : INTEGER
    {
        $$ = &IntegerExpression{Literal: $1.Literal}
    }
    | IDENTIFIER
    {
        $$ = &IdentifierExpression{Literal: $1.Literal}
    }
    | LPAREN expression OPERATOR expression RPAREN
    {
        $$ = &CalculateExpression{Left: $2, Operator: $3.Literal, Right: $4}
    }
    | LT IDENTIFIER arguments GT
    {
        $$ = &CallNumberExpression{Name: $2.Literal, Arguments: $3}
    }

%%
