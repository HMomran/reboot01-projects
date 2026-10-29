# Lem-in

## What it is

A partially copied Go graph project intended to route ants from a start room to an end room and simulate their movement.

## How it works

The parser reads an ant count, room coordinates, `##start`/`##end` markers, and undirected links. It validates room names, coordinates, and link references. `FindAllPaths` uses depth-first search with backtracking to collect simple routes. Conflict helpers check whether candidate routes share intermediate rooms. `main.go` then attempts to select routes and invoke an ant simulation.

## Code map

- `main.go`: parsing, route discovery, selection, and simulation orchestration.
- `Parser/Parser.go`: room/graph types and input validation.
- `Funcs/FindAllPaths.go`: recursive route enumeration.
- `Funcs/Conflict.go`, `Funcs/canAdd.go`: route overlap checks.
- `test/`: input fixtures, including invalid examples.

## Run locally

Run these commands from this project’s folder.

The module specifies Go 1.23.4. The intended command is:

```bash
go run . test/example00.txt
```

This command cannot build from the current copy: the imported `lemin/Simulate` package and the `Funcs.BestPath` implementation are absent.

## Current behavior and limitations

No working simulation can be claimed until the missing code is restored. Enumerating all simple paths can grow rapidly for large graphs. Input starts with a positive integer, uses `name x y` for rooms and `room1-room2` for links, and must identify start/end rooms.

## Try it

Read `test/example00.txt` alongside the parser to follow how its rooms and links become a graph. Restore the missing code before attempting an end-to-end simulation.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
