//go:generate goyacc -o icws94_ygen.go -p "corewar" icws94.y
//go:generate awk -f tokens.awk ../lexer/const.go
package parser
