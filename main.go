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
	fmt.Print(PrintBoard(board))
}
