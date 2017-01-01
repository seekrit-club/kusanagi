package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	InitState()
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if input == "quit\n" || err == io.EOF {
			return
		} else {
			fmt.Print(XboardParse(strings.TrimSpace(input)))
		}
	}
}
