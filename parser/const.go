package parser

import "fmt"

type Opcode int
type OpcodeModifier int
type AddressingMode int

func (o Opcode) String() string {
	for k, v := range opcodes {
		if v == o {
			return k
		}
	}
	return "INVALID"
}

func NewOpcode(s string) (Opcode, error) {
	o, ok := opcodes[s]
	if !ok {
		return -1, fmt.Errorf("wrong opcode %q", s)
	}
	return o, nil
}

func (o OpcodeModifier) String() string {
	for k, v := range opcodeModifiers {
		if v == o {
			return k
		}
	}
	return "INVALID"
}

func NewOpcodeModifier(s string) (OpcodeModifier, error) {
	o, ok := opcodeModifiers[s]
	if !ok {
		return -1, fmt.Errorf("wrong opcode modifier %q", s)
	}
	return o, nil
}

func (a AddressingMode) String() string {
	for k, v := range addressingModes {
		if v == a {
			return k
		}
	}
	return "INVALID"
}

func NewAddressingMode(s string) (AddressingMode, error) {
	a, ok := addressingModes[s]
	if !ok {
		return -1, fmt.Errorf("wrong addressing mode %q", s)
	}
	return a, nil
}

const (
	// instruction opcodes
	DAT Opcode = iota
	MOV
	ADD
	SUB
	MUL
	DIV
	MOD
	JMP
	JMZ
	JMN
	DJN
	CMP
	SLT
	SPL
	ORG
	EQU
	END

	// instruction modifiers
	A OpcodeModifier = iota
	B
	AB
	BA
	F
	X
	I
	NONE

	// addressing modes
	Hash AddressingMode = iota
	Dollar
	At
	Lt
	Gt
)

var opcodes = map[string]Opcode{
	"DAT": DAT,
	"MOV": MOV,
	"ADD": ADD,
	"SUB": SUB,
	"MUL": MUL,
	"DIV": DIV,
	"MOD": MOD,
	"JMP": JMP,
	"JMZ": JMZ,
	"JMN": JMN,
	"DJN": DJN,
	"CMP": CMP,
	"SLT": SLT,
	"SPL": SPL,
	"ORG": ORG,
	"EQU": EQU,
	"END": END,
}

var opcodeModifiers = map[string]OpcodeModifier{
	"A":  A,
	"B":  B,
	"AB": AB,
	"BA": BA,
	"F":  F,
	"X":  X,
	"I":  I,
}

var addressingModes = map[string]AddressingMode{
	"#": Hash,
	"$": Dollar,
	"@": At,
	"<": Lt,
	">": Gt,
}
