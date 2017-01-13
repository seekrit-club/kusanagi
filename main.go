package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strings"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var verbose = flag.Bool("v", false, "verbose output")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	board := InitState()
	reader := bufio.NewReader(os.Stdin)
	var result string
	for {
		input, err := reader.ReadString('\n')
		if input == "quit\n" || err == io.EOF {
			return
		} else {
			board, result = XboardParse(strings.TrimSpace(input),
				board, *verbose)
			fmt.Print(result)
		}
	}
}

func InitState() *Board {
	board := new(Board)
	ClearBoard(board)
	return board
}
