package main

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func XboardParse(line string) string {
	words := strings.Split(line, " ")
	if len(words) == 0 {
		return "\n"
	}
	switch words[0] {
	case "divide":
		if len(words) > 1 {
			res, err := strconv.Atoi(words[1])
			if err == nil {
				start := time.Now()
				nodes := Perft(res, global.board, true)
				elapsed := time.Since(start)
				log.Printf("Divide took %s", elapsed)
				return strconv.FormatUint(nodes, 10) + "\n"
			} else {
				return err.Error()
			}
		}
	case "setboard":
		global.board, _ = Parse(strings.TrimPrefix(line, "setboard "))
	case "usermove":
		if len(words) > 1 {
			move, err := ParseMove(words[1])
			if err == nil {
				MakeMove(global.board, move)
			} else {
				return err.Error()
			}
		}
	case "d":
		return PrintBoard(global.board)
	}
	return "\n"
}
