// YACC grammar for ICWS 94 [0]. See [1] for information on the format.
// 0: https://corewar.co.uk/standards/icws94.htm
// 1: https://pubs.opengroup.org/onlinepubs/9699919799/utilities/yacc.html

// Declarations section

// Stuff between %{%} will be copied over verbatim
%{

package parser

import (
	"log/slog"
	"strconv"

	"github.com/pcolladosoto/corewarg/lexer"
)

type label string

type operation struct {
	opcode int
	modifier int
}

type term struct {
	label label
	immediate int
}

type operand struct {
	mode int
	expr term
}

type instruction struct {
	labels []label
	operation operation
	operandA operand
	operandB operand
}

var programAST []instruction

%}

// Declare the type for values in the stack as well as available
// tag names to declare token and non-terminal types.
%union {
	num int
	label label
	operation operation
	term term
	labelList []label
	comment string
	instruction instruction
	list []instruction
}

%type <operation> operation
%type <num> opcode
%type <num> modifier
%type <num> mode
%type <term> expr
%type <term> term
%type <labelList> label_list
%type <comment> comment
%type <instruction> instruction
%type <instruction> line
%type <list> list
%type <list> assembly_file

// Declare numeral tokens with the same number as declared in the lexer (i.e. lexer.ItemNumber)
// This list is automatically generated with an awk(1) script that's also triggered with 'go generate'.
%token <label> LABEL    1
%token <num> OPCODE     2
%token <num> EOF        3
%token <num> MODIFIER   4
%token <num> MODE       5
%token <num> NUMBER     6
%token <num> OPERAND    7
%token <num> EOL        8
%token <num> KEYWORD    9
%token <num> DAT        10
%token <num> MOV        11
%token <num> ADD        12
%token <num> SUB        13
%token <num> MUL        14
%token <num> DIV        15
%token <num> MOD        16
%token <num> JMP        17
%token <num> JMZ        18
%token <num> JMN        19
%token <num> DJN        20
%token <num> CMP        21
%token <num> SLT        22
%token <num> SPL        23
%token <num> ORG        24
%token <num> EQU        25
%token <num> END        26
%token <num> A  27
%token <num> B  28
%token <num> AB 29
%token <num> BA 30
%token <num> F  31
%token <num> X  32
%token <num> I  33
%token <num> HASH       34
%token <num> DOLLAR     35
%token <num> AT 36
%token <num> LT 37
%token <num> GT 38

// End the declarations
%%

assembly_file:
	list {slog.Warn("redn' at assembly_file", "LIST", $1); $$ = $1; programAST = $1}

list:
	  line      {slog.Warn("redn' at list", "LINE", $1)            ; $$ = []instruction{$1}}
	| line list {slog.Warn("redn' at list", "LINE", $1, "LIST", $2); $$ = append($2, $1)}

line:
	  instruction {slog.Warn("redn' at line", "INSTRUCTION", $1); $$ = $1}
	| comment     {slog.Warn("redn' at line", "COMMENT", $1)}

comment: EOL {slog.Warn("redn' at comment", "EOL", $1)}

instruction:
	  label_list operation mode expr           comment {
		slog.Warn("redn' at instruction", "LABEL_LIST", $1, "OPERATION", $2, "MODE", $3, "EXPR", $4, "COMMENT", $5);
		$$ = instruction{labels: $1, operation: $2, operandA: operand{mode: $3, expr: $4}}
	}
	|            operation mode expr           comment {
		slog.Warn("redn' at instruction", "OPERATION", $1, "MODE", $2, "EXPR", $3, "COMMENT", $4);
		$$ = instruction{labels: nil, operation: $1, operandA: operand{mode: $2, expr: $3}}
	}
	| label_list operation mode expr mode expr comment {
		slog.Warn("redn' at instruction", "LABEL_LIST", $1, "OPERATION", $2, "MODE", $3, "EXPR", $4, "MODE", $5, "EXPR", $6, "COMMENT", $7);
		$$ = instruction{labels: $1, operation: $2, operandA: operand{mode: $3, expr: $4}, operandB: operand{mode: $5, expr: $6}}
	}
	|            operation mode expr mode expr comment {
		slog.Warn("redn' at instruction", "OPERATION", $1, "MODE", $2, "EXPR", $3, "MODE", $4, "EXPR", $5, "COMMENT", $6);
		$$ = instruction{labels: nil, operation: $1, operandA: operand{mode: $2, expr: $3}, operandB: operand{mode: $4, expr: $5}}
	}

label_list:
	  LABEL                {slog.Warn("redn' at label_list", "LABEL", $1)                             ; $$ = []label{$1}}
	| LABEL label_list     {slog.Warn("redn' at label_list", "LABEL", $1, "LABEL_LIST", $2)           ; $$ = append($2, $1)}
	| LABEL EOL label_list {slog.Warn("redn' at label_list", "LABEL", $1, "EOL", $2, "LABEL_LIST", $3); $$ = append($3, $1)}

operation:
	  opcode          {slog.Warn("reduction at operation", "OPCODE", $1)                ; $$ = operation{$1, -1}}
	| opcode modifier {slog.Warn("reduction at operation", "OPCODE", $1, "MODIFIER", $2); $$ = operation{$1, $2}}

opcode:
	  DAT {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| MOV {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| ADD {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| SUB {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| MUL {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| DIV {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| MOD {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| JMP {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| JMZ {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| JMN {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| DJN {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| CMP {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| SLT {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| SPL {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| ORG {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| EQU {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}
	| END {slog.Warn("redn' at opcode", "OPCODE", $1); $$ = $1}

modifier:
	  A  {slog.Warn("redn' at modifier", "MODIFIER", $1); $$ = $1}
	| B  {slog.Warn("redn' at modifier", "MODIFIER", $1); $$ = $1}
	| AB {slog.Warn("redn' at modifier", "MODIFIER", $1); $$ = $1}
	| BA {slog.Warn("redn' at modifier", "MODIFIER", $1); $$ = $1}
	| F  {slog.Warn("redn' at modifier", "MODIFIER", $1); $$ = $1}
	| X  {slog.Warn("redn' at modifier", "MODIFIER", $1); $$ = $1}
	| I  {slog.Warn("redn' at modifier", "MODIFIER", $1); $$ = $1}

mode:
	  HASH        {slog.Warn("redn' at mode", "MODE", $1);      $$ = $1}
	| DOLLAR      {slog.Warn("redn' at mode", "MODE", $1);      $$ = $1}
	| AT          {slog.Warn("redn' at mode", "MODE", $1);      $$ = $1}
	| LT          {slog.Warn("redn' at mode", "MODE", $1);      $$ = $1}
	| GT          {slog.Warn("redn' at mode", "MODE", $1);      $$ = $1}
	| /* empty */ {slog.Warn("redn' at mode", "MODE", "EMPTY"); $$ = -1}

expr:
	term {slog.Warn("reduction at expr", "TERM", $1); $$ = $1}

term:
	  LABEL  {slog.Warn("redn' at term",  "LABEL", $1); $$ = term{label: $1, immediate: -1}}
	| NUMBER {slog.Warn("redn' at term", "NUMBER", $1); $$ = term{immediate: $1, label: ""}}

%%

// This struct should adhere to the corewarLexer interface:
//
//	type corewarLexer interface {
//		Lex(lval *exprSymType) int
//		Error(s string)
//	}
//
// The interface definition is generated by GoYacc!
// Note the prefix (i.e. coreWar) is provided to goyacc
// through the -p flag.
type corewarLex struct {
	l *lexer.Lexer
}

// Lex should return a new token. It's called by the parser. One
// can set the returned token's value through the reference to
// the exprSymType.
func (x *corewarLex) Lex(yylval *corewarSymType) int {
	ni := x.l.NextItem()
	slog.Info("got item", "typ", ni.Typ, "val", ni.Val)

	switch ni.Typ {
	case lexer.ItemNumber:
		pInt, err := strconv.ParseInt(ni.Val, 10, 32)
		if err != nil {
			slog.Error("error parsing number %q: %v\n", ni.Val, err)
		}
		yylval.num = int(pInt)

		return int(lexer.ItemNumber)

	case lexer.ItemLabel:
		yylval.label = label(ni.Val)
		return int(ni.Typ)

	case lexer.ItemEOF:
		return 0 // GoYacc expects EOF to be 0

	default:
		yylval.num = int(ni.Typ)
		return int(ni.Typ)
	}
}

func (x *corewarLex) Error(s string) {
	slog.Error("parse error", "err", s)
}
