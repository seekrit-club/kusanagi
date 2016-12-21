package main

import (
	"testing"
)

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
