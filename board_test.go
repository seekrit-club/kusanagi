package main

import (
	"testing"
)

func TestParseLoadEnPassantNone(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	if OnBoard(board.EnPassant) {
		t.Fail()
	}
}

func TestParseLoadEnPassantA1(t *testing.T) {
	/* This is obviously a bogus position, but we're more interested in
	 * testing the parser here than accurate positions. */
	board, err := Parse("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq a1 0 2")
	if err != nil {
		t.FailNow()
	}
	if board.EnPassant != A1 {
		t.Fail()
	}
}

func TestParseLoadEnPassantH8(t *testing.T) {
	board, err := Parse("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq h8 0 2")
	if err != nil {
		t.FailNow()
	}
	if board.EnPassant != H8 {
		t.Fail()
	}
}

func TestParseToMoveForStart(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	if board.ToMove != WHITE {
		t.Fail()
	}
}

func TestParseToMoveForE4(t *testing.T) {
	board, err := Parse("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 2")
	if err != nil {
		t.FailNow()
	}
	if board.ToMove != BLACK {
		t.Fail()
	}
}

func TestParseDataForSicilian(t *testing.T) {
	board, err := Parse("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2")
	if err != nil {
		t.FailNow()
	}
	if PrintBoard(board) != "rnbqkbnr\npp.ppppp\n........\n..p.....\n....P...\n........\nPPPP.PPP\nRNBQKBNR\n" {
		t.Fail()
	}
}

func TestParseDataForStart(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	if PrintBoard(board) != "rnbqkbnr\npppppppp\n........\n........\n........\n........\nPPPPPPPP\nRNBQKBNR\n" {
		t.Fail()
	}
}

func TestParseInvalidData(t *testing.T) {
	_, err := Parse("rnb%kbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2")
	if err == nil {
		t.Fail()
	}
}

func TestParseInvalidActiveColor(t *testing.T) {
	_, err := Parse("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR o KQkq c6 0 2")
	if err == nil {
		t.Fail()
	}
}

func TestParseInvalidCastling(t *testing.T) {
	_, err := Parse("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w %Qkq c6 0 2")
	if err == nil {
		t.Fail()
	}
}

func TestParseNoCastling(t *testing.T) {
	_, err := Parse("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w - c6 0 2")
	if err != nil {
		t.Fail()
	}
}

func TestParseNoEnPassant(t *testing.T) {
	_, err := Parse("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2")
	if err != nil {
		t.Fail()
	}
}

func TestByteToStringGivenOffboard(t *testing.T) {
	s := ByteToString(OFFBOARD)
	if s != "" {
		t.Fail()
	}
}

func TestByteToStringGivenInvalidPiece(t *testing.T) {
	s := ByteToString(KING + 1)
	if s != "?" {
		t.Fail()
	}
}

func TestIndexToCartesian(t *testing.T) {
	x, y := IndexToCartesian(A1)
	if x != 0 && y != 0 {
		t.Fail()
	}
}

func TestCartesianToAlgebraic(t *testing.T) {
	x := CartesianToAlgebraic(0, 0)
	if x != "a1" {
		t.Fail()
	}
}

func TestIndexToAlgebraic(t *testing.T) {
	x := IndexToAlgebraic(A1)
	if x != "a1" {
		t.Fail()
	}
}

func TestFindKingValid(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	sq, err := FindKing(board, WHITE)
	if err != nil || sq != CartesianToIndex(4, 0) {
		t.Fail()
	}
}

func TestFindKingInvalid(t *testing.T) {
	board, err := Parse("8/8/8/8/8/8/8/8 w - - 0 1")
	if err != nil {
		t.FailNow()
	}
	sq, err := FindKing(board, WHITE)
	if err == nil || sq != INVALID {
		t.Fail()
	}
}

func TestIllegalOnLegalPosition(t *testing.T) {
	board, err := Parse("rnbqkbnr/ppp1pppp/8/1B1p4/4P3/8/PPPP1PPP/RNBQK1NR b KQkq - 1 2")
	if err != nil {
		t.FailNow()
	}
	if Illegal(board) {
		t.Fail()
	}
}

func TestIllegalOnIllegalPosition(t *testing.T) {
	board, err := Parse("rnbqkbnr/ppp1pppp/8/1B1p4/4P3/8/PPPP1PPP/RNBQK1NR w KQkq - 1 2")
	if err != nil {
		t.FailNow()
	}
	if !Illegal(board) {
		t.Fail()
	}
}

func TestIllegalOnEmptyPosition(t *testing.T) {
	board, err := Parse("8/8/8/8/8/8/8/8 w KQkq - 1 2")
	if err != nil {
		t.FailNow()
	}
	if !Illegal(board) {
		t.Fail()
	}
}

func TestNotIllegalWhenSliderBlocked(t *testing.T) {
	board, _ := Parse("rnbq1bnr/pppkpppp/8/3N4/8/8/PPPPPPPP/R1BQKBNR w KQkq - 0 1")
	if Illegal(board) {
		t.Fail()
	}
}

func TestAlgebraicToCartesianLegal(t *testing.T) {
	x, y, err := AlgebraicToCartesian("a1")
	if x != 0 || y != 0 || err != nil {
		t.Fail()
	}
}

func TestAlgebraicToCartesianBadFile(t *testing.T) {
	_, _, err := AlgebraicToCartesian("i1")
	if err == nil {
		t.Fail()
	}
}

func TestAlgebraicToCartesianBadRank(t *testing.T) {
	_, _, err := AlgebraicToCartesian("a9")
	if err == nil {
		t.Fail()
	}
}

func TestAlgebraicToCartesianTooLong(t *testing.T) {
	_, _, err := AlgebraicToCartesian("a1u")
	if err == nil {
		t.Fail()
	}
}

func TestAlgebraicToIndex(t *testing.T) {
	x, err := AlgebraicToIndex("a1")
	if x != A1 || err != nil {
		t.Fail()
	}
}

func TestAlgebraicToIndexInvalid(t *testing.T) {
	_, err := AlgebraicToIndex("z1")
	if err == nil {
		t.Fail()
	}
}
