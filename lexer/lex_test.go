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

	l := lex("lineCommentsTest", test)
	for {
		i := l.nextItem()

		switch i.typ {
		case itemEOF:
			return
		default:
			t.Errorf("got an item other than EOF: type: %s, val: %q", i.typ, i.val)
		}
	}
}

func TestLexSingleInstruction(t *testing.T) {
	tests := []struct {
		in   string
		want []item
	}{
		{"target  DAT.F   #0,     #0", []item{{itemLabel, "target"}, {itemDAT, "DAT"}, {itemF, "F"}, {itemHash, "#"}, {itemNumber, "0"}, {itemHash, "#"}, {itemNumber, "0"}}},
		{"target  DAT.F   #-5,   #15", []item{{itemLabel, "target"}, {itemDAT, "DAT"}, {itemF, "F"}, {itemHash, "#"}, {itemNumber, "-5"}, {itemHash, "#"}, {itemNumber, "15"}}},
		{"ADD.AB  #step,   target", []item{{itemADD, "ADD"}, {itemAB, "AB"}, {itemHash, "#"}, {itemLabel, "step"}, {itemLabel, "target"}}},
		{"MOV.AB  #0,     @target", []item{{itemMOV, "MOV"}, {itemAB, "AB"}, {itemHash, "#"}, {itemNumber, "0"}, {itemAt, "@"}, {itemLabel, "target"}}},
		{"JMP.A    start", []item{{itemJMP, "JMP"}, {itemA, "A"}, {itemLabel, "start"}}},
		{"ORG     start", []item{{itemORG, "ORG"}, {itemLabel, "start"}}},
		{"END", []item{{itemEND, "END"}}},
		{"step    EQU      4", []item{{itemLabel, "step"}, {itemEQU, "EQU"}, {itemNumber, "4"}}},
		{"JMP.A    start ; foo", []item{{itemJMP, "JMP"}, {itemA, "A"}, {itemLabel, "start"}}},
		{"foo fii JMP.A    start ; foo", []item{{itemLabel, "foo"}, {itemLabel, "fii"}, {itemJMP, "JMP"}, {itemA, "A"}, {itemLabel, "start"}}},
		{"foo\nfii JMP.A    start", []item{{itemLabel, "foo"}, {itemLabel, "fii"}, {itemJMP, "JMP"}, {itemA, "A"}, {itemLabel, "start"}}},
		{"\n\t\nfoo\nfii\t JMP.A  \t  start", []item{{itemLabel, "foo"}, {itemLabel, "fii"}, {itemJMP, "JMP"}, {itemA, "A"}, {itemLabel, "start"}}},
	}

	for testI, test := range tests {
		l := lex("singleInstructionTest", test.in)
		j := 0
		for {
			i := l.nextItem()

			slog.Info("got item", "testI", testI, "type", i.typ, "val", i.val)

			if i.typ == itemEOF {
				break
			}

			if j > len(test.want)-1 {
				j++
				continue
			}

			if i.typ != test.want[j].typ || i.val != test.want[j].val {
				t.Errorf("test %d: got type: %s, val: %q; want type: %s, val: %q",
					testI, i.typ, i.val, test.want[j].typ, test.want[j].val)
			}

			j++
		}

		if j != len(test.want) {
			t.Errorf("test %d: got %d items, but expected %d", testI, j, len(test.want))
		}
	}
}
