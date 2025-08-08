//go:generate go tool golang.org/x/tools/cmd/goyacc -o icws94_ygen.go -p "corewar" icws94.y
package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
)

var logger *slog.Logger

func (o Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Opcode   string `json:"opcode"`
		Modifier string `json:"modifier"`
	}{Opcode: o.Opcode.String(), Modifier: o.Modifier.String()})
}

func (o Operand) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Mode string `json:"mode"`
		Expr Term   `json:"expr"`
	}{Mode: o.Mode.String(), Expr: o.Expr})
}

func (i Instruction) String() string {
	buff := bytes.Buffer{}

	// labels
	for j, label := range i.Labels {
		if j == len(i.Labels)-1 {
			buff.WriteString(fmt.Sprintf("%s ", label))
			continue
		}
		buff.WriteString(fmt.Sprintf("%s, ", label))
	}

	// opcode
	buff.WriteString(fmt.Sprintf("%s", i.Operation.Opcode))

	// optional opcode modifier
	if i.Operation.Modifier != OPCODE_MODIFIER_INVALID {
		buff.WriteString(fmt.Sprintf(".%s", i.Operation.Modifier))
	}

	buff.WriteString(" ")

	// operands
	for j, operand := range i.Operands {
		if operand.Mode != ADDRESSING_MODE_INVALID {
			buff.WriteString(fmt.Sprintf("%s", operand.Mode))
		}

		if operand.Expr.Label != "" {
			buff.WriteString(fmt.Sprintf("%s", operand.Expr.Label))
		} else {
			buff.WriteString(fmt.Sprintf("%d", operand.Expr.Immediate))
		}

		if j == 0 && len(i.Operands) == 2 {
			buff.WriteString(", ")
		}
	}

	// comment?

	return buff.String()
}
