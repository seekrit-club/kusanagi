package main

import (
	"math"
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

func FindMove(board *Board) *Move {
	nodecount = 0
	line := new(Line)
	score := AlphaBeta(board, 5, -INFINITY, INFINITY, MATE, line)
	retval := &line.Moves[0]
	retval.Score = score
	return retval
}
