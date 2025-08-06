package parser

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/pcolladosoto/corewarg/lexer"
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

	// Make parsing errors actually useful
	corewarErrorVerbose = true

	// parser verbosity; one of [0, 4]
	corewarDebug = 0
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
	fmt.Printf("programAST: %+v\n", programAST)
}

func TestParserNoMode(t *testing.T) {
	tests := []string{"JMP.A    start\n", "ADD.A  0, target\n"}
	for i, test := range tests {
		if rc := corewarParse(&corewarLex{l: lexer.Lex("parseTest", test)}); rc != 0 {
			t.Errorf("test %d failed", i)
		}
	}
}
