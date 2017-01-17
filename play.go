package main

import (
	"math"
)

var Value [7]int = [7]int{
	0,
	1,
	3,
	3,
	5,
	9,
	0,
}

const INFINITY int = int(math.MaxInt32) // I win!

func Evaluate(board *Board) int {
	mate := Mated(board)
	if mate == -1 {
		return MaterialCount(board)
	} else {
		return mate
	}
}

/* Returns negative infinity if checkmate, zero if stalemate, -1 otherwise. */
func Mated(b *Board) int {
	moves := MoveGen(b)
	if len(moves) > 0 {
		return -1
	}
	var king byte
	if b.ToMove == BLACK {
		king = b.BlackKing
	} else {
		king = b.WhiteKing
	}
	if squareattacked(b, king, b.ToMove^BLACK) {
		return -INFINITY
	} else {
		return 0
	}
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

func AlphaBeta(board *Board, depth, alpha, beta int) int {
	if depth <= 0 {
		return Evaluate(board)
	}
	moves := MoveGen(board)
	for _, move := range moves {
		undo := MakeMove(board, &move)
		val := -AlphaBeta(board, depth-1, -beta, -alpha)
		UnmakeMove(board, &move, undo)
		if val >= beta {
			return beta
		}
		if val > alpha {
			alpha = val
		}
	}
	return alpha
}
