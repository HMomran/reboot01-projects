# ASCII Art

## What it is

A Go command-line program that turns printable ASCII text into eight-row text banners using the bundled `Found.txt` font.

## How it works

`main.go` reads one argument and validates its characters. It converts literal `\n` sequences into line breaks, then renders each input line one banner row at a time. For every character and row, `Rprint` calculates the corresponding line in `Found.txt`, scans the font file, and prints that fragment. Joining fragments horizontally forms a large lettered word.

## Code map

- `main.go`: argument handling and rendering loops.
- `function/Validation.go`: printable ASCII validation.
- `function/Rprint.go`: font-file scanning and output.
- `Found.txt`: glyph data.

## Run locally

Run these commands from this project’s folder.

Install Go 1.23.1 or newer. Keep the working directory here so the relative font path resolves.

```bash
go run . "Hello"
go run . 'Hello\nWorld'
```

## Current behavior and limitations

This version has one font and prints to standard output. Passing no arguments currently panics because `input[0]` is read before argument-count validation. Use printable ASCII; actual newline characters fail validation, while literal `\n` is expanded afterward. The font file is reopened for every fragment.

## Try it

Run the two examples and compare the single-line and multi-line banners.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
