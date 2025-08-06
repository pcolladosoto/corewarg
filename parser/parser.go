//go:generate goyacc -o icws94_ygen.go -p "corewar" icws94.y
package parser

import "log/slog"

var logger *slog.Logger
