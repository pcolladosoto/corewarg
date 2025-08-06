//go:generate goyacc -o icws94_ygen.go -p "corewar" icws94.y
package parser

// Run the command below to extract token constants
//go:generate awk -f tokens.awk ../lexer/const.go
