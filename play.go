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
var abort bool

func Evaluate(board *Board) int {
	phase := calcphase(board)
	opening := MaterialCount(board, false)
	endgame := MaterialCount(board, true)
	score := ((opening * (256 - phase)) + (endgame * phase)) / 256
	if board.ToMove == BLACK {
		score = -score
	}
	return score
}

func calcphase(board *Board) int {
	TotalPhase := 24 // PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2
	phase := TotalPhase
	for _, i := range board.PieceList {
		phase -= PValue[GetPiece(board.Data[i])]
	}
	return (phase*256 + (TotalPhase / 2)) / TotalPhase
}

func MaterialCount(b *Board, endgame bool) int {
	var retval int
	for _, i := range b.PieceList {
		if OnBoard(i) && GetPiece(b.Data[i]) != EMPTY {
			if GetSide(b.Data[i]) == WHITE {
				retval += Value[GetPiece(b.Data[i])] + Pst(GetPiece(b.Data[i]), WHITE, i, endgame)
			} else {
				retval -= (Value[GetPiece(b.Data[i])] + Pst(GetPiece(b.Data[i]), BLACK, i, endgame))
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
	if abort {
		return 0
	}
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
		if abort {
			return 0
		}
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
	if abort {
		return 0
	}
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
		if abort {
			return 0
		}
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
	return (Clock/time.Duration(moves+1) - 20*time.Millisecond)
}

func ThinkingOutput(depth, score int, start time.Time, pv *Line) {
	fmt.Println(depth, score, int64(time.Since(start)/time.Millisecond)/10, nodecount, pv)
}

func SleepThread(board *Board, start time.Time) {
	bedoneby := AllotTime(board)
	fmt.Println("# ", Clock, ": allocated ", bedoneby)
	time.Sleep(bedoneby)
	abort = true
}

func FindMove(board *Board) *Move {
	start := time.Now()
	abort = false
	nodecount = 0
	retval := new(Move)
	for depth := 1; ; depth++ {
		line := new(Line)
		score := AlphaBeta(board, depth, -INFINITY, INFINITY, MATE, line)
		ThinkingOutput(depth, score, start, line)
		if !abort {
			retval = &line.Moves[0]
		} else {
			break
		}

		if depth == 1 {
			Clock -= time.Since(start)
			go SleepThread(board, start)
		}
	}
	return retval
}
