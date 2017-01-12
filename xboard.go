package main

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func XboardParse(line string, verbose bool) string {
	if verbose {
		log.Println(line)
	}
	words := strings.Split(line, " ")
	if len(words) == 0 {
		return "\n"
	}
	switch words[0] {
	case "perft":
		if len(words) > 2 {
			depth, err := strconv.Atoi(words[1])
			if err != nil {
				return err.Error()
			}
			var expected uint64
			expected, err = strconv.ParseUint(words[2],10,64)
			if err != nil {
				return err.Error()
			}
			if Perft(depth, global.board, false) == expected {
				return "SUCCESS\n"
			} else {
				return "FAILURE\n"
			}
		}
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
			move, err := ParseMove(global.board, words[1])
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
