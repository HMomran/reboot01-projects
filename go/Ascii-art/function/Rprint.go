package function

import (
	"bufio"
	"fmt"
	"os"
)

func Rprint(num int) {
	File, err := os.Open("Found.txt")
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer File.Close()
	scanner := bufio.NewScanner(File)
	start := 1
	for scanner.Scan() {
		if start == num {
			fmt.Print(scanner.Text())
		}
		start++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error", err)
		return
	}
}
