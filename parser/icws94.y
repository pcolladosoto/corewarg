// YACC grammar for ICWS 94 [0]. See [1] for information on the format.
// 0: https://corewar.co.uk/standards/icws94.htm
// 1: https://pubs.opengroup.org/onlinepubs/9699919799/utilities/yacc.html

// Declarations section

// Stuff between %{%} will be copied over verbatim
%{

package parser

import (
	"strconv"

	"github.com/pcolladosoto/corewarg/lexer"
)

type Comment string

type Label string

type Operation struct {
	Opcode Opcode
	Modifier OpcodeModifier
}

type Term struct {
	Label Label
	Immediate int
}

type Operand struct {
	Mode AddressingMode
	Expr Term
}

type Instruction struct {
	Labels []Label
	Operation Operation
	Operands []Operand
	Comment Comment
}

var programAST []Instruction

%}

// Declare the type for values in the stack as well as available
// tag names to declare token and non-terminal types.
%union {
	Num int
	Label Label
	Operation Operation
	Term Term
	LabelList []Label
	Comment Comment
	Instruction Instruction
	List []Instruction
	AddressingMode AddressingMode
	Opcode Opcode
	OpcodeModifier OpcodeModifier
}

%type <Operation> operation
%type <AddressingMode> mode
%type <Term> expr
%type <Term> term
%type <LabelList> label_list
%type <Instruction> comment
%type <Instruction> instruction
%type <Instruction> line
%type <List> list
%type <List> assembly_file

// Declare numeral tokens with the same number as declared in the lexer (i.e. lexer.Item*)
// This list is automatically generated with the provided awk(1) script:
//
//   $ awk -f tokens.awk ../lexer/const.go

%token <Num>              EOL             2
%token <Comment>          COMMENT         3
%token <Label>            LABEL           4
%token <Opcode>           OPCODE          5
%token <OpcodeModifier>   OPCODE_MODIFIER 6
%token <AddressingMode>   ADDRESSING_MODE 7
%token <Num>              NUMBER          8
%token <Num>              OPERAND         9

// End the declarations
%%

assembly_file:
	list {
		logger.Debug("redn' at assembly_file", "LIST", $1);
		$$ = $1;

		programAST = $1

		// Reverse the AST as the parser is a bottom up one!
		for i := len(programAST)/2-1; i >= 0; i-- {
			opp := len(programAST)-1-i
			programAST[i], programAST[opp] = programAST[opp], programAST[i]
		}
	}

list:
	  line {
		logger.Debug("redn' at list", "LINE", $1);

		// prevent comments from making it into the AST
		if $1.Operation.Opcode != OPCODE_INVALID {
			$$ = []Instruction{$1}
		} else {
			$$ = nil
		}
	}
	| line list {
		logger.Debug("redn' at list", "LINE", $1, "LIST", $2)

		// prevent comments from making it into the AST
		if $1.Operation.Opcode != OPCODE_INVALID {
			$$ = append($2, $1)
		} else {
			$$ = $2
		}
	}

line:
	  instruction {logger.Debug("redn' at line", "INSTRUCTION", $1); $$ = $1}
	| comment     {logger.Debug("redn' at line", "COMMENT", $1)}

comment:
	  COMMENT EOL {logger.Debug("redn' at comment", "COMMENT", $1, "EOL", $2); $$ = Instruction{Comment: $1}}
	| EOL         {logger.Debug("redn' at comment", "EOL", $1); $$ = Instruction{}}

instruction:
	  label_list operation mode expr           comment {
		logger.Debug("redn' at instruction", "LABEL_LIST", $1, "OPERATION", $2, "MODE", $3, "EXPR", $4, "COMMENT", $5);
		$$ = Instruction{Labels: $1, Operation: $2, Operands: []Operand{{Mode: $3, Expr: $4}}}
	}
	|            operation mode expr           comment {
		logger.Debug("redn' at instruction", "OPERATION", $1, "MODE", $2, "EXPR", $3, "COMMENT", $4);
		$$ = Instruction{Labels: nil, Operation: $1, Operands: []Operand{{Mode: $2, Expr: $3}}}
	}
	| label_list operation mode expr mode expr comment {
		logger.Debug("redn' at instruction", "LABEL_LIST", $1, "OPERATION", $2, "MODE", $3, "EXPR", $4, "MODE", $5, "EXPR", $6, "COMMENT", $7);
		$$ = Instruction{Labels: $1, Operation: $2, Operands: []Operand{{Mode: $3, Expr: $4}, {Mode: $5, Expr: $6}}}
	}
	|            operation mode expr mode expr comment {
		logger.Debug("redn' at instruction", "OPERATION", $1, "MODE", $2, "EXPR", $3, "MODE", $4, "EXPR", $5, "COMMENT", $6);
		$$ = Instruction{Labels: nil, Operation: $1, Operands: []Operand{{Mode: $2, Expr: $3}, {Mode: $4, Expr: $5}}}
	}
	// Special case for END
	| label_list operation comment {
		logger.Debug("redn' at instruction","LABEL_LIST", $1, "OPERATION", $2, "COMMENT", $3);
		$$ = Instruction{Labels: $1, Operation: $2, Operands: nil}
	}
	// Special case for END
	|            operation comment {
		logger.Debug("redn' at instruction", "OPERATION", $1, "COMMENT", $2);
		$$ = Instruction{Labels: nil, Operation: $1, Operands: nil}
	}

label_list:
	  LABEL                {logger.Debug("redn' at label_list", "LABEL", $1)                             ; $$ = []Label{$1}}
	| LABEL label_list     {logger.Debug("redn' at label_list", "LABEL", $1, "LABEL_LIST", $2)           ; $$ = append($2, $1)}
	| LABEL EOL label_list {logger.Debug("redn' at label_list", "LABEL", $1, "EOL", $2, "LABEL_LIST", $3); $$ = append($3, $1)}

operation:
	  OPCODE                 {logger.Debug("redn' at operation", "OPCODE", $1)                       ; $$ = Operation{$1, OPCODE_MODIFIER_INVALID}}
	| OPCODE OPCODE_MODIFIER {logger.Debug("redn' at operation", "OPCODE", $1, "OPCODE_MODIFIER", $2); $$ = Operation{$1, $2}}

mode:
	  ADDRESSING_MODE {logger.Debug("redn' at mode", "ADDRESSING_MODE", $1);      $$ = $1}
	| /* empty */     {logger.Debug("redn' at mode", "ADDRESSING_MODE", "EMPTY"); $$ = ADDRESSING_MODE_INVALID}

expr:
	term {logger.Debug("reduction at expr", "TERM", $1); $$ = $1}

term:
	  LABEL  {logger.Debug("redn' at term",  "LABEL", $1); $$ = Term{Label: $1, Immediate: 0}}
	| NUMBER {logger.Debug("redn' at term", "NUMBER", $1); $$ = Term{Label: "", Immediate: $1}}

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
	logger.Debug("got item", "typ", ni.Typ, "val", ni.Val)

	var err error
	switch ni.Typ {
	case lexer.ItemNumber:
		pInt, err := strconv.ParseInt(ni.Val, 10, 32)
		if err != nil {
			logger.Error("error parsing number %q: %v\n", ni.Val, err)
		}
		yylval.Num = int(pInt)

	case lexer.ItemLabel:
		yylval.Label = Label(ni.Val)

	case lexer.ItemComment:
		yylval.Comment = Comment(ni.Val)

	case lexer.ItemOpcode:
		yylval.Opcode, err = NewOpcode(ni.Val)
		if err != nil {
			logger.Error("error processing opcode", "err", err)
			return -1 // Will this work?
		}

	case lexer.ItemOpcodeModifier:
		yylval.OpcodeModifier, err = NewOpcodeModifier(ni.Val)
		if err != nil {
			logger.Error("error processing opcode modifier", "err", err)
			return -1 // Will this work?
		}

	case lexer.ItemAddressingMode:
		yylval.AddressingMode, err = NewAddressingMode(ni.Val)
		if err != nil {
			logger.Error("error processing addressing mode", "err", err)
			return -1 // Will this work?
		}

	case lexer.ItemEOF:
		return 0 // GoYacc expects EOF to be 0

	default:
		yylval.Num = int(ni.Typ)
	}
	return int(ni.Typ)
}

func (x *corewarLex) Error(s string) {
	logger.Error("parse error", "err", s)
}
