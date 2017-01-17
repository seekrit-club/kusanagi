package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func XboardParse(line string, board *Board, verbose bool) (*Board, string) {
	if verbose {
		log.Println(line)
	}
	words := strings.Split(line, " ")
	if len(words) == 0 {
		return board, "\n"
	}
	switch words[0] {
	case "perft":
		if len(words) > 2 {
			depth, err := strconv.Atoi(words[1])
			if err != nil {
				return board, err.Error()
			}
			var expected uint64
			expected, err = strconv.ParseUint(words[2], 10, 64)
			if err != nil {
				return board, err.Error()
			}
			if Perft(depth, board, false) == expected {
				return board, "SUCCESS\n"
			} else {
				return board, "FAILURE\n"
			}
		}
	case "divide":
		if len(words) > 1 {
			res, err := strconv.Atoi(words[1])
			if err == nil {
				start := time.Now()
				nodes := Perft(res, board, true)
				elapsed := time.Since(start)
				log.Printf("Divide took %s", elapsed)
				return board, strconv.FormatUint(nodes, 10) + "\n"
			} else {
				return board, err.Error()
			}
		}
	case "setboard":
		board, _ = Parse(strings.TrimPrefix(line, "setboard "))
	case "usermove":
		if len(words) > 1 {
			move, err := ParseMove(board, words[1])
			if err == nil {
				MakeMove(board, move)
			} else {
				return board, err.Error()
			}
		}
	case "go":
		move := FindMove(board)
		return board, fmt.Sprintln(move, move.Score)
	case "d":
		return board, PrintBoard(board)
	}
	return board, "\n"
}
