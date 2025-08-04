package lexer

import "fmt"

// Item represents a token or text string returned from the scanner.
type Item struct {
	Typ ItemType
	Val string
}

func (i Item) String() string {
	switch {
	case i.Typ == ItemEOF:
		return "EOF"
	case i.Typ == ItemError:
		return i.Val
	case i.Typ > ItemKeyword:
		return fmt.Sprintf("<%s>", i.Val)
	case len(i.Val) > 10:
		return fmt.Sprintf("%.10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}

// ItemType identifies the type of lex items.
type ItemType int

const (
	ItemError    ItemType = iota // error occurred; value is text of error
	ItemLabel                    // an instruction label
	ItemOpcode                   // an instruction opcode
	ItemEOF                      // EOF
	ItemModifier                 // an instruction modifier
	ItemMode                     // an instruction mode
	ItemNumber                   // an integer number
	ItemOperand                  // a valid operand for an expression

	// simply break up the constants ino two regions to facilitate checks
	ItemKeyword

	// instruction opcodes
	ItemDAT
	ItemMOV
	ItemADD
	ItemSUB
	ItemMUL
	ItemDIV
	ItemMOD
	ItemJMP
	ItemJMZ
	ItemJMN
	ItemDJN
	ItemCMP
	ItemSLT
	ItemSPL
	ItemORG
	ItemEQU
	ItemEND

	// instruction modifiers
	ItemA
	ItemB
	ItemAB
	ItemBA
	ItemF
	ItemX
	ItemI

	// addressing modes
	ItemHash
	ItemDollar
	ItemAt
	ItemLt
	ItemGt
)

// Make the types prettyprint.
var itemName = map[ItemType]string{
	ItemError:    "error",
	ItemLabel:    "label",
	ItemOpcode:   "opcode",
	ItemEOF:      "EOF",
	ItemModifier: "modifier",
	ItemMode:     "mode",
	ItemNumber:   "number",
	ItemOperand:  "operand",

	// opcodes
	ItemDAT: "DAT",
	ItemMOV: "MOV",
	ItemADD: "ADD",
	ItemSUB: "SUB",
	ItemMUL: "MUL",
	ItemDIV: "DIV",
	ItemMOD: "MOD",
	ItemJMP: "JMP",
	ItemJMZ: "JMZ",
	ItemJMN: "JMN",
	ItemDJN: "DJN",
	ItemCMP: "CMP",
	ItemSLT: "SLT",
	ItemSPL: "SPL",
	ItemORG: "ORG",
	ItemEQU: "EQU",
	ItemEND: "END",

	// instruction modifiers
	ItemA:  "A",
	ItemB:  "B",
	ItemAB: "AB",
	ItemBA: "BA",
	ItemF:  "F",
	ItemX:  "X",
	ItemI:  "I",

	// addressing modes
	ItemHash:   "#",
	ItemDollar: "$",
	ItemAt:     "@",
	ItemLt:     "<",
	ItemGt:     ">",
}

var key = map[string]ItemType{
	"DAT": ItemDAT,
	"MOV": ItemMOV,
	"ADD": ItemADD,
	"SUB": ItemSUB,
	"MUL": ItemMUL,
	"DIV": ItemDIV,
	"MOD": ItemMOD,
	"JMP": ItemJMP,
	"JMZ": ItemJMZ,
	"JMN": ItemJMN,
	"DJN": ItemDJN,
	"CMP": ItemCMP,
	"SLT": ItemSLT,
	"SPL": ItemSPL,
	"ORG": ItemORG,
	"EQU": ItemEQU,
	"END": ItemEND,

	"A":  ItemA,
	"B":  ItemB,
	"AB": ItemAB,
	"BA": ItemBA,
	"F":  ItemF,
	"X":  ItemX,
	"I":  ItemI,

	"#": ItemHash,
	"$": ItemDollar,
	"@": ItemAt,
	"<": ItemLt,
	">": ItemGt,
}

func (i ItemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

const eof = -1
