package main

import "testing"

func TestXboardParseSetAndPrintBoard(t *testing.T) {
	InitState()
	XboardParse("setboard " + START)
	board, _ := Parse(START)
	expected := PrintBoard(board)
	if XboardParse("d") != expected {
		t.Log(XboardParse("d"))
		t.Fail()
	}
}
