package main

import (
	"fmt"
	"lemin/Funcs"
	"lemin/Parser"
	simulate "lemin/Simulate"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	timer := time.Now()
	if len(os.Args) != 2 {
		fmt.Println("Invalid format: USAGE go run . <file.txt>")
		return
	}

	file := os.Args[1]

	if !strings.HasSuffix(file, ".txt") {
		fmt.Println("Invalid file: must end with .txt")
		return
	}

	graph, err := Parser.Parser(file)
	if err != nil {
		fmt.Println("ERROR: invalid data format:", err)
		return
	}

	visited := make(map[string]bool)
	var allPaths [][]*Parser.Room

	Funcs.FindAllPaths(graph, visited, graph.Start, []*Parser.Room{}, &allPaths)

	
	if graph.Ants <= 0 {
		fmt.Println("Error: Invalid number of ants")
		return
	}

	fmt.Println("Graph parsed successfully!")
	fmt.Printf("Number of ants: %d\n", graph.Ants)
	fmt.Printf("Start room: %s\n", graph.Start.Name)
	fmt.Printf("End room: %s\n", graph.End.Name)
	fmt.Printf("Total rooms: %d\n", len(graph.Rooms))

	if len(allPaths) == 0 {
		fmt.Println("Error: No paths found from start to end")
		return
	}

	if graph.Ants == 1 {
		sort.Slice(allPaths, func(i, j int) bool {
			return len(allPaths[i]) < len(allPaths[j])
		})
		simulate.SimulateAnts([][]*Parser.Room{allPaths[0]}, graph.Ants)
		fmt.Printf("\nTotal program runtime: %s\n", time.Since(timer))
		return
	}

	fmt.Printf("\nFound %d paths from start to end:\n", len(allPaths))
	for i, path := range allPaths {
		fmt.Printf("Path %d: ", i+1)
		for _, room := range path {
			fmt.Printf("%s ", room.Name)
		}
		fmt.Println()
	}

	best := Funcs.BestPath(allPaths, graph.Start.Name, graph.End.Name)

	fmt.Println("\nBest chosen paths:")
	for i, path := range best {
		fmt.Printf("Path %d: ", i+1)
		for _, r := range path {
			fmt.Print(r.Name, " ")
		}
		fmt.Println()
	}
	fmt.Println("\nSimulation:")
	simulate.SimulateAnts(best, graph.Ants)
	finishTime := time.Since(timer)
	fmt.Printf("\nTotal program runtime: %s\n", finishTime)
}
