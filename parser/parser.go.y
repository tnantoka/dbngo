%{
package parser
%}

%union{
    statements []Statement
    statement Statement
    expression  Expression
    token Token
}

%type<statements> statements

%type<statement> statement command
%type<statement> paper pen line

%type<expression> expression

%token<token> NUMBER LF
%token<token> PAPER PEN LINE

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

statement
    : command LF
    | LF { $$ = nil }

command
    : paper
    | pen
    | line

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

expression
    : NUMBER
    {
        $$ = &NumberExpression{Literal: $1.Literal}
    }
%%

func Parse(yylex yyLexer) int {
	return yyParse(yylex)
}
