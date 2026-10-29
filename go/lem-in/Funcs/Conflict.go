package Funcs

import (
	"lemin/Parser"
)

func Conflict(p1, p2 []*Parser.Room, start, end string) bool {
	visited := make(map[string]bool)

	for _, v := range p1 {
		if v.Name != start && v.Name != end {
			visited[v.Name] = true
		}
	}
	for _, v := range p2 {
		if v.Name != start && v.Name != end && visited[v.Name] {
			return true
		}
	}

	return false
}
