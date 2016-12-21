package main

import (
	"fmt"
)

func main() {
	board, err := Parse(START)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	board.Data[41] = QUEEN
	fmt.Print(PrintBoard(board))
	moves := MoveGen(board)
	fmt.Println(moves)
	MakeMove(board, &moves[0])
	fmt.Print(PrintBoard(board))
	moves = MoveGen(board)
	fmt.Println(moves)
	MakeMove(board, &moves[0])
	fmt.Print(PrintBoard(board))
}
