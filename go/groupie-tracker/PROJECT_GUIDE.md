# Groupie Tracker

## What it is

A server-rendered Go music browser that lists artists and shows band details with concert dates grouped by location.

## How it works

For `/`, the handler fetches artist JSON from the configured Groupie Trackers API, decodes it into Go structs, and renders `backend/indix.html`. `/details?id=...` selects an artist by numeric ID, fetches relation data, matches it by ID, and renders `backend/details.html`. Templates produce the HTML; there is no separate frontend build.

## Code map

- `main.go`: list/detail routes and error handlers.
- `function/FetchArtists.go`: artist/relation structs and JSON fetching.
- `backend/indix.html`: artist list; filename spelling matches the code.
- `backend/details.html`: artist and concert view.

## Run locally

Run these commands from this project’s folder.

Requires Go 1.23 or newer and internet access to the configured API.

```bash
go run .
```

Open `http://localhost:8080`.

## Current behavior and limitations

The source uses `https://groupietrackers.herokuapp.com/api/artists` and `/api/relation`; external availability was not checked. Requests fetch data again instead of using a local database or cache. Missing or unknown IDs redirect home; malformed numeric IDs show a bad-request page.

## Try it

Open the artist list, follow a detail link, and try `/details?id=abc` to inspect invalid-ID handling.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
