package main

import (
	"testing"
)

func TestMateInOne(t *testing.T) {
	board, _ := Parse("5k2/Q7/7N/8/8/K/8/8 w - - 0 1")
	to, _ := AlgebraicToIndex("f7")
	from, _ := AlgebraicToIndex("a7")
	move := FindMove(board)
	if move.To != to || move.From != from {
		t.FailNow()
	}
}
