package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const XBOARDFEATURES string = "feature done=0 usermove=1 setboard=1 myname=\"Kusanagi\" sigterm=0 sigint=0 debug=1 ping=1 colors=0 done=1\n" // our response to the protover command

func XboardParse(line string, board *Board, verbose bool, engine_side *byte) (*Board, string) {
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
	case "new":
		board, _ = Parse(START)
	case "usermove":
		if len(words) > 1 {
			move, err := ParseMove(board, words[1])
			if err == nil {
				MakeMove(board, move)
				board.Moves++
			} else {
				return board, err.Error()
			}
		}
	case "go":
		*engine_side = board.ToMove
                return board, ""
	case "d":
		return board, PrintBoard(board)
	case "protover":
		return board, XBOARDFEATURES
	case "xboard", "post", "nopost", "random":
		return board, ""
	case "ping":
		if len(words) > 1 {
			return board, fmt.Sprintf("pong %s\n", words[1])
		}
	case "time":
		if len(words) > 1 {
			duration, err := time.ParseDuration(words[1] + "0ms")
			if err == nil {
				Clock = duration
			} else {
				return board, err.Error()
			}
		}
	case "level":
		if len(words) > 3 {
			tr, err := strconv.Atoi(words[1])
			if err != nil {
				return board, err.Error()
			}
			var tptc time.Duration
			tptcSpl := strings.Split(words[2], ":")
			if len(tptcSpl) > 1 {
				dur := fmt.Sprintf("%sm%ss", tptcSpl[0],
					tptcSpl[1])
				tptc, err = time.ParseDuration(dur)
				if err != nil {
					return board, err.Error()
				}
			} else {
				tptc, err = time.ParseDuration(words[2] + "m")
				if err != nil {
					return board, err.Error()
				}
			}
			var ti time.Duration
			ti, err = time.ParseDuration(words[3] + "s")
			if err != nil {
				return board, err.Error()
			}
			TimeRepeat = tr
			TimePerTC = tptc
			TimeInc = ti
			Clock = tptc
			return board, fmt.Sprintln("#", tr, tptc, ti)
		}
	}
	return board, ""
}
