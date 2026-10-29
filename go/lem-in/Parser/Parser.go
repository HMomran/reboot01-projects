package Parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	Name  string
	X, Y  int
	Links []*Room
}

type Graph struct {
	Rooms      map[string]*Room
	Start, End *Room
	Ants       int
}

func Parser(file string) (*Graph, error) {
	input, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)

	graph := &Graph{
		Rooms: make(map[string]*Room),
	}

	lineNum := 0
	parsingRooms := true
	var nextIsStart, nextIsEnd bool
	UsedCordinates := make(map[string]bool)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++

		used := false

		if line == "" {
			used = true
			continue
		}

		if lineNum == 1 {
			ants, err := strconv.Atoi(line)
			if err != nil || ants <= 0 {
				return nil, fmt.Errorf("invalid ant count on line 1: %v", err)
			}
			graph.Ants = ants
			used = true
			continue
		}

		if strings.ToLower(line) == "##start" {
			nextIsStart = true
			used = true
			continue
		} else if strings.ToLower(line) == "##end" {
			nextIsEnd = true
			used = true
			continue
		} else if strings.HasPrefix(line, "##") {
			return nil, fmt.Errorf("invalid comment format at line %d: %s", lineNum, line)
		}

		if parsingRooms && strings.Contains(line, " ") {
			parts := strings.Fields(line)
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid room format at line %d: %s", lineNum, line)
			}
			roomName := parts[0]
			if _, exists := graph.Rooms[roomName]; exists {
				return nil, fmt.Errorf("duplicate room name at line %d: %s", lineNum, line)
			} else if roomName[0] == 'L' || roomName[0] == '#' {
				return nil, fmt.Errorf("Room name cannot start with L or # at line %d", lineNum)
			}

			x, err1 := strconv.Atoi(parts[1])
			y, err2 := strconv.Atoi(parts[2])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid coordinates at line %d: %s", lineNum, line)
			}

			Key := fmt.Sprintf("%d,%d", x, y)
			if UsedCordinates[Key] {
				return nil, fmt.Errorf("duplicated coordinates at line %d: %s", lineNum, line)
			}
			UsedCordinates[Key] = true

			room := &Room{
				Name:  parts[0],
				X:     x,
				Y:     y,
				Links: []*Room{},
			}
			graph.Rooms[room.Name] = room

			if nextIsStart {
				graph.Start = room
				nextIsStart = false
			} else if nextIsEnd {
				graph.End = room
				nextIsEnd = false
			}

			used = true
			continue

		} else if strings.Contains(line, "-") {
			parsingRooms = false

			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid link format at line %d: %s", lineNum, line)
			}
			r1, ok1 := graph.Rooms[parts[0]]
			r2, ok2 := graph.Rooms[parts[1]]
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("link references unknown room at line %d: %s", lineNum, line)
			} else if r1 == r2 {
				return nil, fmt.Errorf("a room cannot link to itself at line %d: %s", lineNum, line)
			}
			r1.Links = append(r1.Links, r2)
			r2.Links = append(r2.Links, r1)

			used = true
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}else if !used {
			return nil, fmt.Errorf("invalid or unrecognized line at line %d: %s", lineNum, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	if graph.Start == nil || graph.End == nil {
		return nil, fmt.Errorf("missing start or end room")
	}

	return graph, nil
}
