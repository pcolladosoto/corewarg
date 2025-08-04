package lexer

import "fmt"

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError    itemType = iota // error occurred; value is text of error
	itemLabel                    // an instruction label
	itemOpcode                   // an instruction opcode
	itemEOF                      // EOF
	itemModifier                 // an instruction modifier
	itemMode                     // an instruction mode
	itemNumber                   // an integer number
	itemOperand                  // a valid operand for an expression

	// simply break up the constants ino two regions to facilitate checks
	itemKeyword

	// instruction opcodes
	itemDAT
	itemMOV
	itemADD
	itemSUB
	itemMUL
	itemDIV
	itemMOD
	itemJMP
	itemJMZ
	itemJMN
	itemDJN
	itemCMP
	itemSLT
	itemSPL
	itemORG
	itemEQU
	itemEND

	// instruction modifiers
	itemA
	itemB
	itemAB
	itemBA
	itemF
	itemX
	itemI

	// addressing modes
	itemHash
	itemDollar
	itemAt
	itemLt
	itemGt
)

// Make the types prettyprint.
var itemName = map[itemType]string{
	itemError:    "error",
	itemLabel:    "label",
	itemOpcode:   "opcode",
	itemEOF:      "EOF",
	itemModifier: "modifier",
	itemMode:     "mode",
	itemNumber:   "number",
	itemOperand:  "operand",

	// opcodes
	itemDAT: "DAT",
	itemMOV: "MOV",
	itemADD: "ADD",
	itemSUB: "SUB",
	itemMUL: "MUL",
	itemDIV: "DIV",
	itemMOD: "MOD",
	itemJMP: "JMP",
	itemJMZ: "JMZ",
	itemJMN: "JMN",
	itemDJN: "DJN",
	itemCMP: "CMP",
	itemSLT: "SLT",
	itemSPL: "SPL",
	itemORG: "ORG",
	itemEQU: "EQU",
	itemEND: "END",

	// instruction modifiers
	itemA:  "A",
	itemB:  "B",
	itemAB: "AB",
	itemBA: "BA",
	itemF:  "F",
	itemX:  "X",
	itemI:  "I",

	// addressing modes
	itemHash:   "#",
	itemDollar: "$",
	itemAt:     "@",
	itemLt:     "<",
	itemGt:     ">",
}

var key = map[string]itemType{
	"DAT": itemDAT,
	"MOV": itemMOV,
	"ADD": itemADD,
	"SUB": itemSUB,
	"MUL": itemMUL,
	"DIV": itemDIV,
	"MOD": itemMOD,
	"JMP": itemJMP,
	"JMZ": itemJMZ,
	"JMN": itemJMN,
	"DJN": itemDJN,
	"CMP": itemCMP,
	"SLT": itemSLT,
	"SPL": itemSPL,
	"ORG": itemORG,
	"EQU": itemEQU,
	"END": itemEND,

	"A":  itemA,
	"B":  itemB,
	"AB": itemAB,
	"BA": itemBA,
	"F":  itemF,
	"X":  itemX,
	"I":  itemI,

	"#": itemHash,
	"$": itemDollar,
	"@": itemAt,
	"<": itemLt,
	">": itemGt,
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

const eof = -1
