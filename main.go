package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	InitState()
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		if input == "exit\n" {
			return
		} else {
			fmt.Print(XboardParse(strings.TrimSpace(input)))
		}
	}
}
