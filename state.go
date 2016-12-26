package main

/* Filthy state :^) my program was looking too functional */

type State struct {
	board *Board
}

var global *State

func InitState() {
	global = new(State)
	global.board = new(Board)
	ClearBoard(global.board)
}
