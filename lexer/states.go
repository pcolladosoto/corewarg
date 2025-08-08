package lexer

import (
	"log/slog"
	"strings"
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

const (
	commentDelim = ';'
)

func lexLine(l *Lexer) stateFn {
	slog.Debug("entering lexLine", "start", l.start, "pos", l.pos, "c", string(l.peek()))

	// gobble up leading whitespace
	for isSpace(l.next()) {
		l.ignore()
	}
	l.backup()
	c := l.next()

	// lex a full line comment
	if c == commentDelim {
		return lexComment
	}

	// we're done with the input
	if c == eof {
		l.emit(ItemEOF)
		return nil
	}

	// restore the last rune and lex an instruction
	l.backup()
	return lexInstruction
}

func lexInstruction(l *Lexer) stateFn {
	slog.Debug("entering lexInstruction", "start", l.start, "pos", l.pos, "c", string(l.peek()))
	for {
		switch r := l.next(); {
		case '0' <= r && r <= '9':
			l.backup()
			return lexNumber
		case isAlphaNumeric(r):
			l.backup()
			return lexIdentifier
		case r == '.': // instruction mode
			l.ignore()
			return lexIdentifier
		case strings.Index("#$@<>", string(r)) != -1: // addressing mode
			l.emit(key[string(r)])
			continue
		case strings.Index("+-*/%()", string(r)) != -1: // operand
			l.emit(key[string(r)])
			continue
		case r == commentDelim: // gobble trailing comments
			return lexComment
		case isEOL(r):
			l.emit(ItemEOL)
			return lexLine
		case r == eof:
			// l.ignore() // handle EOLs within a label_list
			return lexLine
		case isSpace(r): // ignore whitespace
			l.ignore()
		}
	}
}

// lexIdentifier scans an alphanumeric or field.
func lexIdentifier(l *Lexer) stateFn {
	slog.Debug("entering lexIdentifier", "start", l.start, "pos", l.pos, "c", string(l.peek()))
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			switch {
			case key[word] > 0:
				l.emit(key[word])
			default:
				l.emit(ItemLabel)
			}
			break Loop
		}
	}
	return lexInstruction
}

// lexNumber scans a decimal number This isn't a perfect number scanner!
func lexNumber(l *Lexer) stateFn {
	slog.Debug("entering lexNumber", "start", l.start, "pos", l.pos, "c", string(l.peek()))
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(ItemNumber)
	return lexInstruction
}

func (l *Lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")

	// Gobble up the number
	l.acceptRun("0123456789")

	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}

	return true
}

func lexComment(l *Lexer) stateFn {
	slog.Debug("entering lexLineComment", "start", l.start, "pos", l.pos, "c", string(l.peek()))

	// drop the comment delimiter
	l.ignore()

	// read until the end of line
	for {
		n := l.next()
		slog.Debug("lexComment", "n", string(n), "l.start", l.start, "l.pos", l.pos, "l.width", l.width)
		if isEOL(n) {
			l.backup() // careful, l.peek changes l.width and doesn't rever it!
			l.emit(ItemComment)
			l.next()
			l.emit(ItemEOL)
			return lexLine
		}
		if n == eof {
			l.backup()
			return lexLine
		}
	}
}
