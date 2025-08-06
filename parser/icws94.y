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

%}

// Declare the type for values in the stack as well as available
// tag names to declare token and non-terminal types.
%union {
	num int
	line string
}

%type <line> instruction

// Declare numeral tokens with the same number as declared in the lexer (i.e. lexer.ItemNumber)
// This list is automatically generated with an awk(1) script that's also triggered with 'go generate'.
%token <num> LABEL      1
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
	list {slog.Warn("reduction at assembly_file with list")}

list:
	  line      {slog.Warn("reduction at list with line")}
	| line list {slog.Warn("reduction at list with line, list")}

line:
	  instruction {slog.Warn("reduction at line with instruction", "instruction", $1)}
	| comment     {slog.Warn("reduction at line with comment")}

comment: EOL {slog.Warn("reduction at comment with EOL")}

instruction:
	  label_list operation mode expr           comment {slog.Warn("reduction at instruction with label_list, operation, mode, comment")}
	|            operation mode expr           comment {slog.Warn("reduction at instruction with label_list, operation, mode, comment")}
	| label_list operation mode expr mode expr comment {slog.Warn("reduction at instruction with label_list, operation, mode, expr, mode, expr, comment")}
	|            operation mode expr mode expr comment {slog.Warn("reduction at instruction with label_list, operation, mode, expr, mode, expr, comment")}

label_list:
	  LABEL                {slog.Warn("reduction at label_list with LABEL")}
	| LABEL label_list     {slog.Warn("reduction at label_list with LABEL, label_list")}
	| LABEL EOL label_list {slog.Warn("reduction at label_list with LABEL, EOL")}

operation:
	  opcode          {slog.Warn("reduction at operation with opcode")}
	| opcode modifier {slog.Warn("reduction at operation with opcode, modifier")}

opcode:
	  DAT {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| MOV {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| ADD {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| SUB {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| MUL {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| DIV {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| MOD {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| JMP {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| JMZ {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| JMN {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| DJN {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| CMP {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| SLT {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| SPL {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| ORG {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| EQU {slog.Warn("reduction at opcode", "OPCODE", $1)}
	| END {slog.Warn("reduction at opcode", "OPCODE", $1)}

modifier:
	A    {slog.Warn("reduction at modifier", "MOD", $1)}
	| B  {slog.Warn("reduction at modifier", "MOD", $1)}
	| AB {slog.Warn("reduction at modifier", "MOD", $1)}
	| BA {slog.Warn("reduction at modifier", "MOD", $1)}
	| F  {slog.Warn("reduction at modifier", "MOD", $1)}
	| X  {slog.Warn("reduction at modifier", "MOD", $1)}
	| I  {slog.Warn("reduction at modifier", "MOD", $1)}

mode:
	HASH          {slog.Warn("reduction at mode", "MODE", $1)}
	| DOLLAR      {slog.Warn("reduction at mode", "MODE", $1)}
	| AT          {slog.Warn("reduction at mode", "MODE", $1)}
	| LT          {slog.Warn("reduction at mode", "MODE", $1)}
	| GT          {slog.Warn("reduction at mode", "MODE", $1)}
	| /* empty */ {slog.Warn("reduction at mode", "MODE", "EMPTY")}

expr:
	term {slog.Warn("reduction at expr with term")}

term:
	  LABEL  {slog.Warn("reduction at term with LABEL")}
	| NUMBER {slog.Warn("reduction at term with NUMBER")}

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

	case lexer.ItemEOF:
		return 0 // GoYacc expects EOF to be 0

	default:
		return int(ni.Typ)
	}
}

func (x *corewarLex) Error(s string) {
	slog.Error("parse error", "err", s)
}
