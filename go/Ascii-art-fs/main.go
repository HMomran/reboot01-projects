package main

import (
	"asciiart/function"
	"fmt"
	"os"
)

func main() {
	var result string
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . [STRING] [BANNER]" + "\n" + "EX: go run . something standard")
		return
	}
	input := os.Args[1:]
	word := input[0]
	banner := input[1]

	if !(function.Validation(word)) {
		fmt.Println("Invalid input")
		os.Exit(0)
	}

	result = function.AsciiArt(word, banner)

	fmt.Print(result)

}
