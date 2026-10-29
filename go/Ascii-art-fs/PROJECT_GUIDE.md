# ASCII Art — Font Selection

## What it is

A Go banner renderer that lets the caller choose standard, shadow, or thinkertoy lettering.

## How it works

`main.go` validates arguments and passes the text and font name to `AsciiArt`. The renderer expands literal `\n` separators and joins eight rows of glyph fragments for each line. The font reader first tries local files and falls back to a hardcoded GitHub raw URL when opening a local file fails. The completed string is printed to the terminal.

## Code map

- `main.go`: command-line interface.
- `function/AsciiArt.go`: line and glyph assembly.
- `function/PrintAscii.go`: font lookup and HTTP fallback.
- `function/Validation.go`: ASCII character checks.
- `Standard.txt`, `Shadow.txt`, `Thinkertoy.txt`: supplied font files.

## Run locally

Run these commands from this project’s folder.

Requires Go 1.23.1 or newer.

```bash
go run . "Hello" standard
go run . 'Hello\nWorld' shadow
```

## Current behavior and limitations

The code looks for `banner/standard.txt`, `banner/shadow.txt`, and `banner/ThinkerToy.txt`, but this copy stores the supplied fonts at the project root with different capitalization. As copied, the HTTP fallback is therefore used. For offline use, create `banner/` and copy the supplied fonts to those exact expected filenames. The fallback can make many requests because it is performed per glyph row. Font names are normalized to lowercase. This project selects fonts; it does not save output files itself.

## Try it

Render a short word with each font. Try a literal `\n` between two words.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
