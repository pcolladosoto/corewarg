package lexer

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time.
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}

			// Remove the directory from the source's filename.
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}
			return a
		},
	}))
	slog.SetDefault(logger)
}

func TestLexLineComments(t *testing.T) {
	test := `
	; this is a comment
	;this is a comment too!
	; empty lines should be okay too!

	`

	l := Lex("lineCommentsTest", test)
	for {
		i := l.NextItem()

		switch i.Typ {
		case ItemEOF:
			return
		default:
			t.Errorf("got an item other than EOF: type: %s, val: %q", i.Typ, i.Val)
		}
	}
}

func TestLexSingleInstruction(t *testing.T) {
	tests := []struct {
		in   string
		want []Item
	}{
		{"target  DAT.F   #0,     #0", []Item{{ItemLabel, "target"}, {ItemOpcode, "DAT"}, {ItemOpcodeModifier, "F"}, {ItemAddressingMode, "#"}, {ItemNumber, "0"}, {ItemAddressingMode, "#"}, {ItemNumber, "0"}}},
		{"target  DAT.F   #-5,   #15", []Item{{ItemLabel, "target"}, {ItemOpcode, "DAT"}, {ItemOpcodeModifier, "F"}, {ItemAddressingMode, "#"}, {ItemNumber, "-5"}, {ItemAddressingMode, "#"}, {ItemNumber, "15"}}},
		{"ADD.AB  #step,   target", []Item{{ItemOpcode, "ADD"}, {ItemOpcodeModifier, "AB"}, {ItemAddressingMode, "#"}, {ItemLabel, "step"}, {ItemLabel, "target"}}},
		{"MOV.AB  #0,     @target", []Item{{ItemOpcode, "MOV"}, {ItemOpcodeModifier, "AB"}, {ItemAddressingMode, "#"}, {ItemNumber, "0"}, {ItemAddressingMode, "@"}, {ItemLabel, "target"}}},
		{"JMP.A    start", []Item{{ItemOpcode, "JMP"}, {ItemOpcodeModifier, "A"}, {ItemLabel, "start"}}},
		{"ORG     start", []Item{{ItemOpcode, "ORG"}, {ItemLabel, "start"}}},
		{"END", []Item{{ItemOpcode, "END"}}},
		{"step    EQU      4", []Item{{ItemLabel, "step"}, {ItemOpcode, "EQU"}, {ItemNumber, "4"}}},
		{"JMP.A    start ; foo", []Item{{ItemOpcode, "JMP"}, {ItemOpcodeModifier, "A"}, {ItemLabel, "start"}}},
		{"foo fii JMP.A    start ; foo", []Item{{ItemLabel, "foo"}, {ItemLabel, "fii"}, {ItemOpcode, "JMP"}, {ItemOpcodeModifier, "A"}, {ItemLabel, "start"}}},
		{"foo\nfii JMP.A    start", []Item{{ItemLabel, "foo"}, {ItemEOL, "\n"}, {ItemLabel, "fii"}, {ItemOpcode, "JMP"}, {ItemOpcodeModifier, "A"}, {ItemLabel, "start"}}},
		{"\n\t\nfoo\nfii\t JMP.A  \t  start", []Item{{ItemLabel, "foo"}, {ItemEOL, "\n"}, {ItemLabel, "fii"}, {ItemOpcode, "JMP"}, {ItemOpcodeModifier, "A"}, {ItemLabel, "start"}}},
	}

	for testI, test := range tests {
		l := Lex("singleInstructionTest", test.in)
		j := 0
		for {
			i := l.NextItem()

			slog.Info("got item", "testI", testI, "type", i.Typ, "val", i.Val)

			if i.Typ == ItemEOF {
				break
			}

			if j > len(test.want)-1 {
				j++
				continue
			}

			if i.Typ != test.want[j].Typ || i.Val != test.want[j].Val {
				t.Errorf("test %d: got type: %s, val: %q; want type: %s, val: %q",
					testI, i.Typ, i.Val, test.want[j].Typ, test.want[j].Val)
			}

			j++
		}

		if j != len(test.want) {
			t.Errorf("test %d: got %d items, but expected %d", testI, j, len(test.want))
		}
	}
}
