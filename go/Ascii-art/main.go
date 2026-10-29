package main

import (
	"asciiart/function"
	"fmt"
	"os"
	"strings"
)

func main() {
	input := os.Args[1:]
	word := input[0]

	if len(input) != 1 {
		fmt.Println("error: empty input")
		return
	}

	if !function.Validation(word) {
		return
	}

	processed := strings.ReplaceAll(word, `\n`, "\n")
	lines := strings.Split(processed, "\n")

	for j := 0; j < len(lines); j++ {
		if lines[j] == "" && j == len(lines)-1 {
			break
		}else if lines[j] == "" {
			fmt.Println()
			continue
		}

		for k := 0; k < 8; k++ {
			for m := 0; m < len(lines[j]); m++ {
				function.Rprint(((int(lines[j][m])-31)*9 - 7) + k)
			}
			fmt.Println()
		}
	}
}
