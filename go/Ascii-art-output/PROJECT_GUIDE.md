# ASCII Art — File Output

## What it is

A Go banner renderer that writes ASCII art into a chosen text file.

## How it works

`main.go` validates arguments and passes the text and font name to `AsciiArt`. The renderer expands literal `\n` separators and joins eight rows of glyph fragments for each line. The font reader first tries local files and falls back to a hardcoded GitHub raw URL when opening a local file fails. The resulting string is written using `os.Create`, replacing an existing output file.

## Code map

- `main.go`: command-line interface and file writing.
- `function/AsciiArt.go`: line and glyph assembly.
- `function/AsciiPrint.go`: font lookup and HTTP fallback.
- `function/Validation.go`: ASCII character checks.
- `Standard.txt`, `Shadow.txt`, `Thinkertoy.txt`: supplied font files.

## Run locally

Run these commands from this project’s folder.

Requires Go 1.23.1 or newer.

```bash
go run . --output=result.txt "Hello" standard
cat result.txt
```

## Current behavior and limitations

The code looks for `banner/standard.txt`, `banner/shadow.txt`, and `banner/ThinkerToy.txt`, but this copy stores the supplied fonts at the project root with different capitalization. As copied, the HTTP fallback is therefore used. For offline use, create `banner/` and copy the supplied fonts to those exact expected filenames. The fallback can make many requests because it is performed per glyph row. Font names are normalized to lowercase. Use a `.txt` output name. Output creation overwrites an existing file.

## Try it

Render a short word with each font. Inspect the generated text file in a monospaced editor.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
