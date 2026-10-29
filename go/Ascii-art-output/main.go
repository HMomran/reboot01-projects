package main

import (
	"asciiart/function"
	"fmt"
	"os"
	"strings"
)

func main() {
	var result string
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]" + "\n" + "EX: go run . --output=<fileName.txt> something standard")
		return
	}
	input := os.Args[1:]
	file := strings.ReplaceAll(input[0], " ", "")
	if !strings.HasPrefix(file, "--output=") {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]" + "\n" + "EX: go run . --output=<fileName.txt> something standard")
		return
	}
	word := input[1]
	banner := input[2]

	if !(function.Validation(word)) {
		fmt.Println("Error: Invalid ascii-art input")
		return
	}

	result = function.AsciiArt(word, banner)
	var file2 string
	array := []rune(file)
	for i := 0; i < len(array); i++ {
		if i > 8 {
			file2 += string(array[i])
		}
	}

	var Type string
	for j := len(array) - 1; j > len(array)-4; j-- {
		Type += string(array[j])
	}
	if Type != "txt" {
		fmt.Println("Only files of type .txt")
		os.Exit(0)
	}

	OutPutFile, err := os.Create(file2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer OutPutFile.Close()

	_, err = OutPutFile.WriteString(result)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Done!")

}
