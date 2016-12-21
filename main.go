package main

import (
	"fmt"
)

func main() {
	board := new(Board)
	InitBoard(board)
	fmt.Println(PrintBoard(board))
}
