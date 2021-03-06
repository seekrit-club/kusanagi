package main

import (
	"testing"
)

func TestMakeMovePawnPush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	to := byte(CartesianToIndex(0, 2))
	from := byte(CartesianToIndex(0, 1))
	move := Move{from, to, MoveQuiet, EMPTY, 0}
	MakeMove(board, &move)
	res, err := FindPiece(board, to)
	if GetPiece(board.Data[from]) != EMPTY || GetPiece(board.Data[to]) != PAWN ||
		err != nil || board.PieceList[res] != to {
		t.Fail()
	}
}

func TestMakeMoveWhitePawnDoublePush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	to := byte(CartesianToIndex(0, 3))
	from := byte(CartesianToIndex(0, 1))
	move := Move{from, to, MoveDoublePush, EMPTY, 0}
	MakeMove(board, &move)
	res, err := FindPiece(board, to)
	if GetPiece(board.Data[from]) != EMPTY || GetPiece(board.Data[to]) !=
		PAWN || board.EnPassant != CartesianToIndex(0, 2) ||
		err != nil || board.PieceList[res] != to {
		t.Fail()
	}
}

func TestMakeMoveBlackPawnDoublePush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	board.ToMove = BLACK
	to := byte(CartesianToIndex(0, 4))
	from := byte(CartesianToIndex(0, 6))
	move := Move{from, to, MoveDoublePush, EMPTY, 0}
	MakeMove(board, &move)
	res, err := FindPiece(board, to)
	if GetPiece(board.Data[from]) != EMPTY || GetPiece(board.Data[to]) !=
		PAWN || board.EnPassant != CartesianToIndex(0, 5) ||
		err != nil || board.PieceList[res] != to {
		t.Fail()
	}
}

func TestMakeMoveWhitePawnCapture(t *testing.T) {
	board, err := Parse("8/8/8/3p4/4P3/8/8/8 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	to := byte(CartesianToIndex(3, 4))
	from := byte(CartesianToIndex(4, 3))
	move := Move{from, to, MoveCapture, EMPTY, 0}
	MakeMove(board, &move)
	if GetPiece(board.Data[from]) != EMPTY || GetPiece(board.Data[to]) !=
		PAWN || GetSide(board.Data[to]) != WHITE {
		t.Fail()
	}
}

func TestMakeMoveBlackPawnCapture(t *testing.T) {
	board, err := Parse("8/8/8/3p4/4P3/8/8/8 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	to := byte(CartesianToIndex(4, 3))
	from := byte(CartesianToIndex(3, 4))
	move := Move{from, to, MoveCapture, EMPTY, 0}
	MakeMove(board, &move)
	if GetPiece(board.Data[from]) != EMPTY || GetPiece(board.Data[to]) !=
		PAWN || GetSide(board.Data[to]) != BLACK {
		t.Fail()
	}
}

func IsMoveInMoveList(t *testing.T, list []Move, from, to, kind byte) bool {
	for _, m := range list {
		if m.To == to && m.From == from && m.Kind == kind {
			return true
		}
	}
	return false
}

func TestMoveGenWhitePawnPush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(0, 1)), byte(CartesianToIndex(0, 2)), MoveQuiet) {
		t.Fail()
	}
}

func TestMoveGenBlackPawnPush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	board.ToMove = BLACK
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(0, 6)), byte(CartesianToIndex(0, 5)), MoveQuiet) {
		t.Fail()
	}
}

func TestMoveGenWhitePawnDoublePush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(0, 1)), byte(CartesianToIndex(0, 3)), MoveDoublePush) {
		t.Fail()
	}
}

func TestMoveGenBlackPawnDoublePush(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	board.ToMove = BLACK
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(0, 6)), byte(CartesianToIndex(0, 4)), MoveDoublePush) {
		t.Fail()
	}
}

func TestMoveGenSlider(t *testing.T) {
	board, err := Parse("4k3/8/8/1r6/8/8/8/4K3 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(1, 4)), byte(CartesianToIndex(1, 7)), MoveQuiet) {
		t.Fail()
	}
}

func TestMoveGenNonSlider(t *testing.T) {
	board, err := Parse("4k3/8/8/1r6/8/8/8/4K3 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(4, 7)), byte(CartesianToIndex(4, 6)), MoveQuiet) {
		t.Fail()
	}
	if IsMoveInMoveList(t, moves, byte(CartesianToIndex(4, 7)), byte(CartesianToIndex(4, 3)), MoveQuiet) {
		t.Fail()
	}
}

func TestMoveGenWhitePawnCap(t *testing.T) {
	board, err := Parse("4k3/8/8/3p4/4P3/8/8/4K3 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(4, 3)), byte(CartesianToIndex(3, 4)), MoveCapture) {
		t.Fail()
	}
}

func TestMoveGenBlackPawnCap(t *testing.T) {
	board, err := Parse("4k3/8/8/3p4/4P3/8/8/4K3 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(3, 4)), byte(CartesianToIndex(4, 3)), MoveCapture) {
		t.Fail()
	}
}

func TestMoveGenSliderCap(t *testing.T) {
	board, err := Parse("4k3/8/8/8/1r1P4/8/8/4K3 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(1, 3)), byte(CartesianToIndex(3, 3)), MoveCapture) {
		t.Fail()
	}
}

func TestMoveGenNonSliderCap(t *testing.T) {
	board, err := Parse("4k3/8/2n5/8/3P4/8/8/4K3 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	t.Log(moves)
	if !IsMoveInMoveList(t, moves, byte(CartesianToIndex(2, 5)), byte(CartesianToIndex(3, 3)), MoveCapture) {
		t.Fail()
	}
}

func TestPerftStart(t *testing.T) {
	tperft(t, 1, 20)
	tperft(t, 2, 400)
	tperft(t, 3, 8902)
	tperft(t, 4, 197281)
	tperft(t, 5, 4865609)
}

func TestPerftKiwipete(t *testing.T) {
	board, _ :=
		Parse("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -")
	tperftboard(t, 1, 48, board)
	tperftboard(t, 2, 2039, board)
	tperftboard(t, 3, 97862, board)
	tperftboard(t, 4, 4085603, board)
}

func tperftboard(t *testing.T, depth int, expected uint64, b *Board) {
	result := Perft(depth, b, false)
	t.Log("Result of perft(", depth, "): ", result)
	if result != expected {
		t.FailNow()
	}
}

func tperft(t *testing.T, depth int, expected uint64) {
	result := DoDividePerft(depth)
	t.Log("Result of perft(", depth, "): ", result)
	if result != expected {
		t.FailNow()
	}
}

func TestSquareNotAttacked(t *testing.T) {
	board, err := Parse("4k3/8/8/8/8/8/3p4/4K3 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	if squareattacked(board, CartesianToIndex(4, 4), BLACK) {
		t.FailNow()
	}
}

func TestSquareAttackedByWhitePawn(t *testing.T) {
	board, err := Parse("4k3/3P4/8/8/8/8/8/4K3 b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	if !squareattacked(board, CartesianToIndex(4, 7), WHITE) {
		t.FailNow()
	}
}

func TestSquareAttackedByBlackPawn(t *testing.T) {
	board, err := Parse("4k3/8/8/8/8/8/3p4/4K3 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	if !squareattacked(board, CartesianToIndex(4, 0), BLACK) {
		t.FailNow()
	}
}

func TestSquareAttackedByQueenFile(t *testing.T) {
	board, err := Parse("4k3/8/8/8/8/4q3/8/4K3 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	if !squareattacked(board, CartesianToIndex(4, 0), BLACK) {
		t.FailNow()
	}
}

func TestSquareAttackedByQueenDiag(t *testing.T) {
	board, err := Parse("4k3/8/8/8/1q6/8/8/4K3 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	if !squareattacked(board, CartesianToIndex(4, 0), BLACK) {
		t.FailNow()
	}
}

func TestSquareAttackedByRook(t *testing.T) {
	board, err := Parse("4k3/8/8/8/8/4r3/8/4K3 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	if !squareattacked(board, CartesianToIndex(4, 0), BLACK) {
		t.FailNow()
	}
}

func TestSquareAttackedByBishop(t *testing.T) {
	board, err := Parse("4k3/8/8/8/1b6/8/8/4K3 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	if !squareattacked(board, CartesianToIndex(4, 0), BLACK) {
		t.FailNow()
	}
}

func TestSquareAttackedByKnight(t *testing.T) {
	board, err := Parse("4k3/8/8/8/8/8/2n5/4K3 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	if !squareattacked(board, CartesianToIndex(4, 0), BLACK) {
		t.FailNow()
	}
}

func TestMakeMoveEnPassantWhite(t *testing.T) {
	board, err := Parse("rnbqkbnr/1pp1pppp/p7/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 1")
	if err != nil {
		t.FailNow()
	}
	from, _ := AlgebraicToIndex("e5")
	to, _ := AlgebraicToIndex("d6")
	taken, _ := AlgebraicToIndex("d5")
	MakeMove(board, &Move{from, to, MoveEnPassant, EMPTY, 0})
	if GetPiece(board.Data[taken]) != EMPTY || GetPiece(board.Data[to]) != PAWN {
		t.Fail()
	}
}

func TestMakeMoveEnPassantBlack(t *testing.T) {
	board, err := Parse("rnbqkbnr/ppp1pppp/8/8/3pP3/PP6/2PP1PPP/RNBQKBNR b KQkq e3 0 1")
	if err != nil {
		t.FailNow()
	}
	from, _ := AlgebraicToIndex("d4")
	to, _ := AlgebraicToIndex("e3")
	taken, _ := AlgebraicToIndex("e4")
	MakeMove(board, &Move{from, to, MoveEnPassant, EMPTY, 0})
	if GetPiece(board.Data[taken]) != EMPTY || GetPiece(board.Data[to]) != PAWN {
		t.Fail()
	}
}

func TestMoveGenEnPassant(t *testing.T) {
	board, err := Parse("rnbqkbnr/1pp1pppp/p7/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 1")
	if err != nil {
		t.FailNow()
	}
	moves := MoveGen(board)
	from, _ := AlgebraicToIndex("e5")
	to, _ := AlgebraicToIndex("d6")
	if !IsMoveInMoveList(t, moves, from, to, MoveEnPassant) {
		t.Fail()
	}
}

func TestUnmakeMoveBlackEnPassant(t *testing.T) {
	board, err := Parse("rnbqkbnr/ppp1pppp/8/8/3pP3/PP6/2PP1PPP/RNBQKBNR b KQkq e3 0 1")
	if err != nil {
		t.FailNow()
	}
	from, _ := AlgebraicToIndex("d4")
	to, _ := AlgebraicToIndex("e3")
	taken, _ := AlgebraicToIndex("e4")
	move := &Move{from, to, MoveEnPassant, EMPTY, 0}
	undo := MakeMove(board, move)
	UnmakeMove(board, move, undo)
	if GetPiece(board.Data[taken]) != PAWN || GetPiece(board.Data[to]) != EMPTY || GetPiece(board.Data[from]) != PAWN {
		t.Fail()
	}
}

func TestPerftPawnPromotions(t *testing.T) {
	board, err := Parse("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1")
	if err != nil {
		t.FailNow()
	}
	if Perft(4, board, false) != 182838 {
		t.FailNow()
	}
}

func TestMoveToLongAlgebraicPromoteQueen(t *testing.T) {
	from, _ := AlgebraicToIndex("e7")
	to, _ := AlgebraicToIndex("e8")
	move := &Move{from, to, MovePromote, QUEEN, 0}
	if MoveToLongAlgebraic(move) != "e7e8q" {
		t.FailNow()
	}
}

func TestMoveToLongAlgebraicPromoteRook(t *testing.T) {
	from, _ := AlgebraicToIndex("e7")
	to, _ := AlgebraicToIndex("e8")
	move := &Move{from, to, MovePromote, ROOK, 0}
	if MoveToLongAlgebraic(move) != "e7e8r" {
		t.FailNow()
	}
}

func TestMoveToLongAlgebraicPromoteBishop(t *testing.T) {
	from, _ := AlgebraicToIndex("e7")
	to, _ := AlgebraicToIndex("e8")
	move := &Move{from, to, MovePromote, BISHOP, 0}
	if MoveToLongAlgebraic(move) != "e7e8b" {
		t.FailNow()
	}
}

func TestMoveToLongAlgebraicPromoteKnight(t *testing.T) {
	from, _ := AlgebraicToIndex("e7")
	to, _ := AlgebraicToIndex("e8")
	move := &Move{from, to, MovePromote, KNIGHT, 0}
	if MoveToLongAlgebraic(move) != "e7e8n" {
		t.FailNow()
	}
}

func TestParseMoveValid(t *testing.T) {
	board, _ := Parse(START)
	_, err := ParseMove(board, "e2e4")
	if err != nil {
		t.FailNow()
	}
}

func TestParseMoveInvalid(t *testing.T) {
	board, _ := Parse(START)
	_, err := ParseMove(board, "e1e2")
	if err == nil {
		t.FailNow()
	}
}

func TestFilterCapturesOnEmptyList(t *testing.T) {
	board, _ := Parse(START)
	movelist := FilterCaptures(MoveGen(board))
	if len(movelist) != 0 {
		t.FailNow()
	}
}

func TestFilterCapturesOnMixedList(t *testing.T) {
	board, _ := Parse("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2")
	movelist := FilterCaptures(MoveGen(board))
	if len(movelist) != 1 {
		t.FailNow()
	}
}
