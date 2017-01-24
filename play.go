package main

import (
	"fmt"
	"math"
	"time"
)

var Value [7]int = [7]int{
	// Value of the pieces in centipawns
	0,
	100,
	300,
	300,
	500,
	900,
	0,
}

type Line struct {
	Moves []Move
}

const INFINITY int = int(math.MaxInt32) // I win!
const MATE int = INFINITY - 10          // Value of a checkmate in 1

var nodecount uint64

func Evaluate(board *Board) int {
	return MaterialCount(board)
}

func MaterialCount(b *Board) int {
	var retval int
	for i := A1; i <= H8; i++ {
		if OnBoard(i) && GetPiece(b.Data[i]) != EMPTY {
			if GetSide(b.Data[i]) == b.ToMove {
				retval += Value[GetPiece(b.Data[i])]
			} else {
				retval -= Value[GetPiece(b.Data[i])]
			}
		}
	}
	return retval
}

func AlphaBeta(board *Board, depth, alpha, beta, mate int, pline *Line) int {
	nodecount++
	legal := 0
	if depth <= 0 {
		return Evaluate(board)
	}
	line := new(Line)
	moves := MoveGen(board)
	for _, move := range moves {
		undo := MakeMove(board, &move)
		if Illegal(board) {
			UnmakeMove(board, &move, undo)
			continue
		}
		val := -AlphaBeta(board, depth-1, -beta, -alpha, mate-1, line)
		UnmakeMove(board, &move, undo)
		if val >= beta {
			return beta
		}
		if val > alpha {
			alpha = val
			pline.Moves = append([]Move{move}, pline.Moves...)
			line.Moves = append(pline.Moves, move)
		}
		legal++
	}
	if legal == 0 {
		if !InCheck(board) {
			return 0
		}
		return -mate
	}
	return alpha
}

func AllotTime(board *Board) time.Duration {
	repeat := TimeRepeat
	if repeat <= 0 {
		repeat = 40
	}
	moves := repeat - board.Moves
	for moves <= 0 {
		moves += repeat
	}
	return (Clock/time.Duration(moves+1))/10 - 20*time.Millisecond
}

func ThinkingOutput(depth, score int, start time.Time, pv *Line) {
	fmt.Println(depth, score, int64(time.Since(start) / time.Millisecond)/10, nodecount, pv)
}

func FindMove(board *Board) *Move {
	start := time.Now()
	timetomove := AllotTime(board)
	bedoneby := start.Add(timetomove)
	nodecount = 0
	for depth := 1; ; depth++ {
		line := new(Line)
		score := AlphaBeta(board, depth, -INFINITY, INFINITY, MATE, line)
		ThinkingOutput(depth, score, start, line)
		retval := &line.Moves[0]
		if time.Now().After(bedoneby) {
			retval.Score = score
			Clock -= time.Since(start)
			return retval
		}
	}
}
