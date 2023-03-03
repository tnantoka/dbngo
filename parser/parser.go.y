%{
package parser
%}

%union{
    statements []Statement
    statement Statement
    expression Expression
    token Token
}

%type<statements> statements body

%type<statement> statement command
%type<statement> paper pen line set dot copy repeat
%type<statement> block

%type<expression> expression

%token<token> NUMBER LF IDENTIFIER OPERATOR
%token<token> PAPER PEN LINE SET REPEAT
%token<token> LBRACE RBRACE LPAREN RPAREN LBRACKET RBRACKET

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

expression
    : NUMBER
    {
        $$ = &NumberExpression{Literal: $1.Literal}
    }
    | IDENTIFIER
    {
        $$ = &IdentifierExpression{Literal: $1.Literal}
    }
    | LPAREN expression OPERATOR expression RPAREN
    {
        $$ = &CalculateExpression{Left: $2, Operator: $3.Literal, Right: $4}
    }
%%
