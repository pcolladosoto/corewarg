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
		want []itemType
	}{
		{"target  DAT.F   #0,     #0", []itemType{itemLabel, itemDAT, itemF, itemHash, itemNumber, itemHash, itemNumber}},
		{"target  DAT.F   #-5,   #15", []itemType{itemLabel, itemDAT, itemF, itemHash, itemNumber, itemHash, itemNumber}},
		{"ADD.AB  #step,   target", []itemType{itemADD, itemAB, itemHash, itemLabel, itemLabel}},
		{"MOV.AB  #0,     @target", []itemType{itemMOV, itemAB, itemHash, itemNumber, itemAt, itemLabel}},
		{"JMP.A    start", []itemType{itemJMP, itemA, itemLabel}},
		{"ORG     start", []itemType{itemORG, itemLabel}},
		{"END", []itemType{itemEND}},
		{"step    EQU      4", []itemType{itemLabel, itemEQU, itemNumber}},
		{"JMP.A    start ; foo", []itemType{itemJMP, itemA, itemLabel}},
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

			if i.typ != test.want[j] {
				t.Errorf("test %d: got type: %s, val: %q; want type: %s", testI, i.typ, i.val, test.want[j])
			}

			j++
		}

		if j != len(test.want) {
			t.Errorf("test %d: got %d items, but expected %d", testI, j, len(test.want))
		}
	}
}
