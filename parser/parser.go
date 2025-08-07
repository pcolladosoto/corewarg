//go:generate goyacc -o icws94_ygen.go -p "corewar" icws94.y
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

var foo = Instruction{}

var rawTemplate string = `
{{if .Labels}}
{{range $i, $label := .Labels}}{{ $label }}, 
{{end}}
{{end}}

{{if ne .Operation.Opcode 0 }}

This release contains the changes made between tags {{.FromTag}} and {{.ToTag}}.

You can [compare it with the previous release](https://github.com/scitags/flowd-go/compare/{{.FromTag}}...{{.ToTag}}).

{{if .Feat}}
{{if .GhMarkdown}}
<details>
<summary><h2>Features</h2></summary>
{{else}}
## Features
{{end}}
{{range $hash, $description := .Feat}}1. {{ $description }} ({{ $hash }})
{{end}}
{{if .GhMarkdown}}
</details>
{{end}}
{{end}}

{{if .Fixes}}
{{if .GhMarkdown}}
<details>
<summary><h2>Fixes</h2></summary>
{{else}}
## Fixes
{{end}}
{{range $hash, $description := .Fixes}}1. {{ $description }} ({{ $hash }})
{{end}}
{{if .GhMarkdown}}
</details>
{{end}}
{{end}}

{{if .Ci}}
{{if .GhMarkdown}}
<details>
<summary><h2>CI</h2></summary>
{{else}}
## CI
{{end}}
{{range $hash, $description := .Ci}}1. {{ $description }} ({{ $hash }})
{{end}}
{{if .GhMarkdown}}
</details>
{{end}}
{{end}}

{{if .Docs}}
{{if .GhMarkdown}}
<details>
<summary><h2>Docs</h2></summary>
{{else}}
## Docs
{{end}}
{{range $hash, $description := .Docs}}1. {{ $description }} ({{ $hash }})
{{end}}
{{if .GhMarkdown}}
</details>
{{end}}
{{end}}

{{if .Build}}
{{if .GhMarkdown}}
<details>
<summary><h2>Build</h2></summary>
{{else}}
## Build
{{end}}
{{range $hash, $description := .Build}}1. {{ $description }} ({{ $hash }})
{{end}}
{{if .GhMarkdown}}
</details>
{{end}}
{{end}}

{{if .Perf}}
{{if .GhMarkdown}}
<details>
<summary><h2>Performance</h2></summary>
{{else}}
## Performance
{{end}}
{{range $hash, $description := .Perf}}1. {{ $description }} ({{ $hash }})
{{end}}
{{if .GhMarkdown}}
</details>
{{end}}
{{end}}

{{if .Test}}
{{if .GhMarkdown}}
<details>
<summary><h2>Tests</h2></summary>
{{else}}
## Tests
{{end}}
{{range $hash, $description := .Test}}1. {{ $description }} ({{ $hash }})
{{end}}
{{if .GhMarkdown}}
</details>
{{end}}
{{end}}

{{if .Authors}}
{{ $addEmail := .AddEmail}}
## Contributors
{{range $email, $name := .Authors}}1. **{{$name}}**{{if $addEmail}} <{{ $email }}>{{end}}
{{end}}
{{end}}

`

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
