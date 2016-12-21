package main

import (
	"testing"
)

func TestParseForSicilian(t *testing.T) {
	board, err := Parse("rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2")
	if err != nil {
		t.FailNow()
	}
	if PrintBoard(board) != "rnbqkbnr\npp.ppppp\n........\n..p.....\n....P...\n........\nPPPP.PPP\nRNBQKBNR\n" {
			t.Fail()
		}
}

func TestParseForStart(t *testing.T) {
	board, err := Parse(START)
	if err != nil {
		t.FailNow()
	}
	if PrintBoard(board) != "rnbqkbnr\npppppppp\n........\n........\n........\n........\nPPPPPPPP\nRNBQKBNR\n" {
			t.Fail()
		}
}
