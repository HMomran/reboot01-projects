# Math Skills

## What it is

A Go utility that reads one number per line from a text file and prints the average, median, population variance, and standard deviation.

## How it works

`main.go` scans the file and parses each line as `float64`. `Average` divides the sum by the count. `Median` sorts the slice and selects the middle value or averages the two middle values. `Variance` averages squared distances from the mean using population size N; `StandardDeviation` takes its square root. The program rounds each final result to an integer.

## Code map

- `main.go`: input parsing and rounded output.
- `functions/`: individual statistical calculations.
- `input.txt`: example input.

## Run locally

Run these commands from this project’s folder.

Requires Go 1.23.1 or newer.

```bash
go run . input.txt
```

## Current behavior and limitations

Use a nonempty `.txt` file with one valid number on every line, without blank lines. Empty input is not safely handled. This calculates population variance, not sample variance. Median sorts the input slice in place. Short filenames can trigger an indexing error in the extension check.

## Try it

Create a text file containing 1, 2, and 3 on separate lines. Average and median should both print 2; population variance rounds to 1 and standard deviation rounds to 1.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
