# ASCII Art Web

## What it is

A Go website where a visitor submits text and selects a banner font to generate ASCII art.

## How it works

The server renders `templates/index.html` at `/`. A POST to `/result` reads the `Data` and `banner` form fields, calls the banner renderer, and inserts the output into `templates/result.html`. CSS is served under `/templates/`. The renderer validates printable ASCII and a maximum of 50 input bytes, expands literal `\n`, and assembles eight-row glyphs. Error templates handle invalid input and rendering failures.

## Code map

- `main.go`: routes, form handling, and template rendering.
- `functions/`: validation, glyph lookup, and banner assembly.
- `templates/`: form, result page, CSS, and error pages.
- `banner/`: supplied font data.

## Run locally

Run these commands from this project’s folder.

Requires Go 1.23.1 or newer.

```bash
go run .
```

Open `http://localhost:8080`.

## Current behavior and limitations

The code requests lowercase `banner/standard.txt` and `banner/shadow.txt`, plus `banner/ThinkerToy.txt`; supplied names use different capitalization. These differ on case-sensitive filesystems, so the HTTP font fallback may be used, including in Linux containers. Rename or copy fonts to the exact expected paths for reliable offline use. Several projects use port 8080; run one at a time. Actual newline characters fail validation; literal `\n` is expanded after validation.

## Try it

Submit a short word with each font, try a literal `\n`, and submit empty or over-50-character text to see validation behavior.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
