package lexer

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
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

type tests []struct {
	in   string
	want []Item
}

func runTests(t *testing.T, ts tests) {
	for i, test := range ts {
		slog.Debug("new test", "i", i)
		l := Lex("lexTest", test.in)
		j := 0
		for {
			item := l.NextItem()

			if item.Typ == ItemEOF {
				break
			}

			if item.Typ != test.want[j].Typ || item.Val != test.want[j].Val {
				t.Errorf("test %d: got type: %s, val: %q; want type: %s, val: %q",
					i, item.Typ, item.Val, test.want[j].Typ, test.want[j].Val,
				)
				break
			}

			j++
		}
		if j != len(test.want) {
			t.Errorf("test %d: got %d items, but expected %d", i, j, len(test.want))
		}
	}
}

func TestLexLineComments(t *testing.T) {
	ts := tests{
		{"", nil},
		{"\n", nil},
		{"\n\n", nil},
		{"; this is a comment", nil}, // we need "\n" terminations
		{"; this is a comment\n", []Item{{ItemComment, " this is a comment"}, {ItemEOL, "\n"}}},
		{";this is a comment too\n", []Item{{ItemComment, "this is a comment too"}, {ItemEOL, "\n"}}},
	}

	runTests(t, ts)
}

func TestLexSingleInstruction(t *testing.T) {
	ts := tests{
		{"target  DAT.F   #0,     #0", []Item{{ItemLabel, "target"}, {ItemOpcode, "DAT"}, {ItemOpcodeModifier, "F"}, {ItemAddressingMode, "#"}, {ItemNumber, "0"}, {ItemAddressingMode, "#"}, {ItemNumber, "0"}}},
		{"target  DAT.F   #-5,   #15", []Item{{ItemLabel, "target"}, {ItemOpcode, "DAT"}, {ItemOpcodeModifier, "F"}, {ItemAddressingMode, "#"}, {ItemOperand, "-"}, {ItemNumber, "5"}, {ItemAddressingMode, "#"}, {ItemNumber, "15"}}},
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

	runTests(t, ts)
}

func TestLexTwoInstructions(t *testing.T) {
	ts := tests{
		{"target  DAT.F   #0,     #0", []Item{{ItemLabel, "target"}, {ItemOpcode, "DAT"}, {ItemOpcodeModifier, "F"}, {ItemAddressingMode, "#"}, {ItemNumber, "0"}, {ItemAddressingMode, "#"}, {ItemNumber, "0"}}},
		{"target  DAT.F   #-5,   #15", []Item{{ItemLabel, "target"}, {ItemOpcode, "DAT"}, {ItemOpcodeModifier, "F"}, {ItemAddressingMode, "#"}, {ItemOperand, "-"}, {ItemNumber, "5"}, {ItemAddressingMode, "#"}, {ItemNumber, "15"}}},
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

	runTests(t, ts)
}

func TestLexFiles(t *testing.T) {
	dataDir := "testdata"
	paths, err := os.ReadDir(dataDir)
	if err != nil {
		t.Fatalf("error reading dir contents: %v", err)
	}

	wants := map[string][]Item{
		"dwarf.rc": {
			{ItemComment, "redcode"},
			{ItemEOL, "\n"},

			{ItemComment, "name          Dwarf"},
			{ItemEOL, "\n"},

			{ItemComment, "author        A. K. Dewdney"},
			{ItemEOL, "\n"},

			{ItemComment, "version       94.1"},
			{ItemEOL, "\n"},

			{ItemComment, "date          April 29, 1993"},
			{ItemEOL, "\n"},

			{ItemComment, "strategy      Bombs every fourth instruction."},
			{ItemEOL, "\n"},

			{ItemOpcode, "ORG"},
			{ItemLabel, "start"},
			{ItemComment, " Indicates the instruction with"},
			{ItemEOL, "\n"},

			{ItemComment, " the label \"start\" should be the"},
			{ItemEOL, "\n"},

			{ItemComment, " first to execute."},
			{ItemEOL, "\n"},

			{ItemLabel, "step"},
			{ItemOpcode, "EQU"},
			{ItemNumber, "4"},
			{ItemComment, " Replaces all occurrences of \"step\""},
			{ItemEOL, "\n"},

			{ItemComment, " with the character \"4\"."},
			{ItemEOL, "\n"},

			{ItemLabel, "target"},
			{ItemOpcode, "DAT"},
			{ItemOpcodeModifier, "F"},
			{ItemAddressingMode, "#"},
			{ItemNumber, "0"},
			{ItemAddressingMode, "#"},
			{ItemNumber, "0"},
			{ItemComment, " Pointer to target instruction."},
			{ItemEOL, "\n"},

			{ItemLabel, "start"},
			{ItemOpcode, "ADD"},
			{ItemOpcodeModifier, "AB"},
			{ItemAddressingMode, "#"},
			{ItemLabel, "step"},
			{ItemLabel, "target"},
			{ItemComment, " Increments pointer by step."},
			{ItemEOL, "\n"},

			{ItemOpcode, "MOV"},
			{ItemOpcodeModifier, "AB"},
			{ItemAddressingMode, "#"},
			{ItemNumber, "0"},
			{ItemAddressingMode, "@"},
			{ItemLabel, "target"},
			{ItemComment, " Bombs target instruction."},
			{ItemEOL, "\n"},

			{ItemOpcode, "JMP"},
			{ItemOpcodeModifier, "A"},
			{ItemLabel, "start"},
			{ItemComment, " Same as JMP.A -2.  Loops back to"},
			{ItemEOL, "\n"},

			{ItemComment, " the instruction labelled \"start\"."},
			{ItemEOL, "\n"},

			{ItemOpcode, "END"},
			{ItemEOL, "\n"},
		},
	}

	ts := tests{}

	for _, file := range paths {
		fName := file.Name()

		prog, err := os.ReadFile(fmt.Sprintf("%s/%s", dataDir, fName))
		if err != nil {
			t.Errorf("error reading file %q: %v", fName, err)
			continue
		}

		want, ok := wants[fName]
		if !ok {
			t.Errorf("no wants defined for file %q", fName)
			continue
		}

		ts = append(ts, struct {
			in   string
			want []Item
		}{in: string(prog), want: want})
	}

	runTests(t, ts)
}

func TestLexOperands(t *testing.T) {
	ts := tests{
		{"a + b\n", []Item{{ItemLabel, "a"}, {ItemOperand, "+"}, {ItemLabel, "b"}, {ItemEOL, "\n"}}},
		{"a + b + c\n", []Item{{ItemLabel, "a"}, {ItemOperand, "+"}, {ItemLabel, "b"}, {ItemOperand, "+"}, {ItemLabel, "c"}, {ItemEOL, "\n"}}},
		{"a - b\n", []Item{{ItemLabel, "a"}, {ItemOperand, "-"}, {ItemLabel, "b"}, {ItemEOL, "\n"}}},
		{"a - b - c\n", []Item{{ItemLabel, "a"}, {ItemOperand, "-"}, {ItemLabel, "b"}, {ItemOperand, "-"}, {ItemLabel, "c"}, {ItemEOL, "\n"}}},
		{"a * b\n", []Item{{ItemLabel, "a"}, {ItemOperand, "*"}, {ItemLabel, "b"}, {ItemEOL, "\n"}}},
		{"a * b * c\n", []Item{{ItemLabel, "a"}, {ItemOperand, "*"}, {ItemLabel, "b"}, {ItemOperand, "*"}, {ItemLabel, "c"}, {ItemEOL, "\n"}}},
		{"a / b\n", []Item{{ItemLabel, "a"}, {ItemOperand, "/"}, {ItemLabel, "b"}, {ItemEOL, "\n"}}},
		{"a / b / c\n", []Item{{ItemLabel, "a"}, {ItemOperand, "/"}, {ItemLabel, "b"}, {ItemOperand, "/"}, {ItemLabel, "c"}, {ItemEOL, "\n"}}},
		{"a % b\n", []Item{{ItemLabel, "a"}, {ItemOperand, "%"}, {ItemLabel, "b"}, {ItemEOL, "\n"}}},
		{"a % b % c\n", []Item{{ItemLabel, "a"}, {ItemOperand, "%"}, {ItemLabel, "b"}, {ItemOperand, "%"}, {ItemLabel, "c"}, {ItemEOL, "\n"}}},
		{"a + (b)\n", []Item{{ItemLabel, "a"}, {ItemOperand, "+"}, {ItemOperand, "("}, {ItemLabel, "b"}, {ItemOperand, ")"}, {ItemEOL, "\n"}}},
		{"a * (b + c)\n", []Item{
			{ItemLabel, "a"}, {ItemOperand, "*"}, {ItemOperand, "("}, {ItemLabel, "b"}, {ItemOperand, "+"}, {ItemLabel, "c"},
			{ItemOperand, ")"}, {ItemEOL, "\n"}},
		},
	}

	runTests(t, ts)
}
