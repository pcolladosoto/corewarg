package parser

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/pcolladosoto/corewarg/lexer"
)

func init() {
	l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
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
	logger = l

	// Make parsing errors actually useful
	corewarErrorVerbose = true

	// parser verbosity; one of [0, 4]
	corewarDebug = 0
}

func marshalAST(ast []Instruction) {
	enc, err := json.MarshalIndent(ast, "", "    ")
	if err != nil {
		fmt.Printf("error marshalling AST: %v\n", err)
		return
	}
	fmt.Printf("programAST: %s\n", enc)
}

func printAST(ast []Instruction) {
	// escape the '%' so that it makes it through...
	format := fmt.Sprintf("%%0%dd", int(math.Ceil(math.Log10(float64(len(ast))))))
	for i, ins := range ast {
		// build the format string before supplying data!
		fmt.Printf(fmt.Sprintf("%s: %%s\n", format), i, ins)
	}
}

// Remember to run 'gp generate' beforehand!
func TestParserSimple(t *testing.T) {
	tests := []string{"ADD.A  #0, #target\n", "ADD.A  #0, #target\nADD.A  #0, #target\n"}
	for i, test := range tests {
		if rc := corewarParse(&corewarLex{l: lexer.Lex("parseTest", test)}); rc != 0 {
			t.Errorf("test %d failed", i)
		}
	}
}

func TestParserError(t *testing.T) {
	tests := []string{"5 WRONG 4\n", "5 WRONG 4"}
	for i, test := range tests {
		if rc := corewarParse(&corewarLex{l: lexer.Lex("parseTest", test)}); rc == 0 {
			t.Errorf("test %d passed and it shouldn't...", i)
		}
	}
}

func TestParserSingleFieldInstruction(t *testing.T) {
	tests := []string{"foo faa JMP.A    #start\n"}
	for i, test := range tests {
		if rc := corewarParse(&corewarLex{l: lexer.Lex("parseTest", test)}); rc != 0 {
			t.Errorf("test %d failed", i)
		}
	}
	marshalAST(programAST)
}

func TestParserNoMode(t *testing.T) {
	tests := []string{"JMP.A    start\n", "ADD.A  0, target\n"}
	for i, test := range tests {
		if rc := corewarParse(&corewarLex{l: lexer.Lex("parseTest", test)}); rc != 0 {
			t.Errorf("test %d failed", i)
		}
	}
	marshalAST(programAST)
}

func TestParserFiles(t *testing.T) {
	dataDir := "../lexer/testdata"
	paths, err := os.ReadDir(dataDir)
	if err != nil {
		t.Fatalf("error reading dir contents: %v", err)
	}

	for i, file := range paths {
		prog, err := os.ReadFile(fmt.Sprintf("%s/%s", dataDir, file.Name()))
		if err != nil {
			t.Errorf("error reading file %q: %v", file.Name(), err)
		}

		if rc := corewarParse(&corewarLex{l: lexer.Lex("parseTest", string(prog))}); rc != 0 {
			t.Errorf("test %d failed", i)
		}
		marshalAST(programAST)
		printAST(programAST)
	}
}
