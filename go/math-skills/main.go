package main

import (
	"bufio"
	"fmt"
	"math"
	"maths/functions"
	"os"
	"strconv"
)

func main() {

	var array []float64

	if len(os.Args) != 2 {
		fmt.Println("Invalid input. You must run the program using: go run . <filename.txt>")
		os.Exit(0)
	}

	input := os.Args[1:]
	reader, err := os.Open(input[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	var Type string
	Array := []rune(input[0])
	for j := len(Array) - 1; j > len(Array)-4; j-- {
		Type += string(Array[j])
	}
	if Type != "txt" {
		fmt.Println("Only files of type .txt are valid")
		os.Exit(0)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		num, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		array = append(array, num)
	}

	fmt.Println("Average: ", int(math.Round(functions.Average(array))))
	fmt.Println("Median: ", int(math.Round(functions.Median(array))))
	fmt.Println("Variance: ", int(math.Round(functions.Variance(array))))
	fmt.Println("StandardDeviation: ", int(math.Round(functions.StandardDeviation(array))))

}
