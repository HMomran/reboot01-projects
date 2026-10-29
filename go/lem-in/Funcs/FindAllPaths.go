package Funcs

import "lemin/Parser"

func FindAllPaths(graph *Parser.Graph, visited map[string]bool, current *Parser.Room, path []*Parser.Room, allPaths *[][]*Parser.Room) {
	visited[current.Name] = true
	path = append(path, current)

	if current == graph.End {
		pathCopy := make([]*Parser.Room, len(path))
		copy(pathCopy, path)
		*allPaths = append(*allPaths, pathCopy)
	} else {
		for _, neighbor := range current.Links {
			if !visited[neighbor.Name] {
				FindAllPaths(graph, visited, neighbor, path, allPaths)
			}
		}
	}

	path = path[:len(path)-1]
	visited[current.Name] = false
}
