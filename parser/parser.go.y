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

%type<statement> statement
%type<statement> command
%type<statement> paper

%type<expression> expression

%token<token> NUMBER LF
%token<token> PAPER

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

paper
    : PAPER expression
    {
        $$ = &PaperStatement{Value: $2}
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
