/*
Package lexer implements a lexer based on Rob Pike's talk [0] for the ICWS 94 Standard [1].
The language grammar is defined as:

```

	assembly_file:
		list

	list:
		line | line list

	line:
		comment | instruction

	comment:
		; v* EOL | EOL

	instruction:
		label_list operation mode field comment |
		label_list operation mode expr , mode expr comment

	label_list:
		label | label label_list | label newline label_list | e

	label:
		alpha alphanumeral*

	operation:
		opcode | opcode.modifier

	opcode:
		DAT | MOV | ADD | SUB | MUL | DIV | MOD |
		JMP | JMZ | JMN | DJN | CMP | SLT | SPL |
		ORG | EQU | END

	modifier:
		A | B | AB | BA | F | X | I

	mode:
		# | $ | @ | < | > | e

	expr:
		term |
		term + expr | term - expr |
		term * expr | term / expr |
		term % expr

	term:
		label | number | (expression)

	number:
		whole_number | signed_integer

	signed_integer:
		+whole_number | -whole_number

	whole_number:
		numeral+

	alpha:
		A-Z | a-z | _

	numeral:
		0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9

	alphanumeral:
		alpha | numeral

	v:
		^EOL

	EOL:
		newline | EOF

	newline:
		LF | CR | LF CR | CR LF

	e:

```

Where 'e' is the "empty" element, meaning the token may be omitted, a
caret '^' means NOT, an asterisk '*' immediately adjacent means zero or
more occurrences of the previous token, and a plus '+' immediately
adjacent means one or more occurrences of the previous token. The
vertical bar '|' means OR.

This lexer is intended to be used together with the goyacc-based parser
provided by the accompanying parser package.

0: https://go.dev/talks/2011/lex/r59-lex.go
1: https://corewar.co.uk/standards/icws94.htm
*/
package lexer
