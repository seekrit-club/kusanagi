package main

import (
	"fmt"
	"math"
	"time"
)

type Line struct {
	Moves []Move
}

const INFINITY int = int(math.MaxInt32) // I win!
const MATE int = INFINITY - 10          // Value of a checkmate in 1

var nodecount uint64

func Evaluate(board *Board) int {
	return MaterialCount(board, false)
}

func MaterialCount(b *Board, endgame bool) int {
	var retval int
	for _, i := range b.PieceList {
		if OnBoard(i) && GetPiece(b.Data[i]) != EMPTY {
			if GetSide(b.Data[i]) == b.ToMove {
				retval += Value[GetPiece(b.Data[i])] + Pst(GetPiece(b.Data[i]), b.ToMove, i, endgame)
			} else {
				retval -= (Value[GetPiece(b.Data[i])] + Pst(GetPiece(b.Data[i]), b.ToMove, i, endgame))
			}
		}
	}
	return retval
}

func Pst(piece, side, index byte, endgame bool) int {
	switch piece {
	case PAWN:
		return GetPst(index, side, pstPawnMg, pstPawnEg, endgame)
	case KNIGHT:
		return GetPst(index, side, pstKnightMg, pstKnightEg, endgame)
	case BISHOP:
		return GetPst(index, side, pstBishopMg, pstBishopEg, endgame)
	case ROOK:
		return GetPst(index, side, pstRookMg, pstRookEg, endgame)
	case QUEEN:
		return GetPst(index, side, pstQueenMg, pstQueenMg, endgame)
	case KING:
		return GetPst(index, side, pstKingMg, pstKingEg, endgame)
	default:
		return 0
	}
}

func GetPst(index, side byte, tableMg [64]int, tableEg [64]int, endgame bool) int {
	var table []int
	if endgame {
		table = tableEg[:]
	} else {
		table = tableMg[:]
	}
	file, rank := IndexToCartesian(index)
	var tableindex byte
	if side == BLACK {
		tableindex = (8 * (7 - rank)) + file
	} else {
		tableindex = (8 * rank) + file
	}
	return table[tableindex]
}

func Quies(board *Board, alpha, beta int) int {
	nodecount++
	eval := Evaluate(board)
	if eval >= beta {
		return beta
	}
	if eval > alpha {
		alpha = eval
	}
	moves := FilterCaptures(MoveGen(board))
	for _, move := range moves {
		undo := MakeMove(board, &move)
		if Illegal(board) {
			UnmakeMove(board, &move, undo)
			continue
		}
		val := -Quies(board, -beta, -alpha)
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

func AlphaBeta(board *Board, depth, alpha, beta, mate int, pline *Line) int {
	nodecount++
	legal := 0
	if depth <= 0 {
		return Quies(board, alpha, beta)
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
	fmt.Println(depth, score, int64(time.Since(start)/time.Millisecond)/10, nodecount, pv)
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
