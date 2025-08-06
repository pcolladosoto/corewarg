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
	case len(i.Val) > 10:
		return fmt.Sprintf("%.10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}

// ItemType identifies the type of lex items.
type ItemType int

const (
	ItemError          ItemType = iota // error occurred; value is text of error
	ItemEOF                            // End Of File (EOF)
	ItemEOL                            // a valid End Of Line (EOL)
	ItemLabel                          // an instruction label (i.e. an alphanumeric)
	ItemOpcode                         // an instruction opcode (i.e. DAT, MOV, ADD, ...)
	ItemOpcodeModifier                 // an instruction modifier (i.e. A, B, AB, BA, F, X, I)
	ItemAddressingMode                 // an instruction addressing mode (i.e. #, $, @, <, >)
	ItemNumber                         // an integer number
	ItemOperand                        // a valid operand for an expression (i.e. +, -, *, /, %)
)

// Make the types prettyprint.
var itemName = map[ItemType]string{
	ItemError:          "error",
	ItemLabel:          "label",
	ItemOpcode:         "opcode",
	ItemEOF:            "EOF",
	ItemOpcodeModifier: "modifier",
	ItemAddressingMode: "mode",
	ItemNumber:         "number",
	ItemOperand:        "operand",
	ItemEOL:            "EOL",
}

var key = map[string]ItemType{
	"DAT": ItemOpcode,
	"MOV": ItemOpcode,
	"ADD": ItemOpcode,
	"SUB": ItemOpcode,
	"MUL": ItemOpcode,
	"DIV": ItemOpcode,
	"MOD": ItemOpcode,
	"JMP": ItemOpcode,
	"JMZ": ItemOpcode,
	"JMN": ItemOpcode,
	"DJN": ItemOpcode,
	"CMP": ItemOpcode,
	"SLT": ItemOpcode,
	"SPL": ItemOpcode,
	"ORG": ItemOpcode,
	"EQU": ItemOpcode,
	"END": ItemOpcode,

	"A":  ItemOpcodeModifier,
	"B":  ItemOpcodeModifier,
	"AB": ItemOpcodeModifier,
	"BA": ItemOpcodeModifier,
	"F":  ItemOpcodeModifier,
	"X":  ItemOpcodeModifier,
	"I":  ItemOpcodeModifier,

	"#": ItemAddressingMode,
	"$": ItemAddressingMode,
	"@": ItemAddressingMode,
	"<": ItemAddressingMode,
	">": ItemAddressingMode,
}

func (i ItemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

const eof = -1
