package Funcs

import "lemin/Parser"

func canAdd(allPaths [][]*Parser.Room, NewPath []*Parser.Room, start, end string) bool {
	for _, v := range allPaths {
		if Conflict(v, NewPath, start, end) {
			return false
		}
	}
	return true
}
