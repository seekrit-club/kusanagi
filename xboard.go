package main

import (
	"strings"
)

func XboardParse(line string) string {
	words := strings.Split(line, " ")
	if len(words) == 0 {
		return "\n"
	}
	switch words[0] {
	case "setboard":
		global.board, _ = Parse(strings.TrimPrefix(line, "setboard "))
	case "d":
		return PrintBoard(global.board)
	}
	return "\n"
}
