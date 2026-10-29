# Forum — Go Templates

## What it is

A Go and SQLite discussion forum with registration, sign-in, categorized posts, comments, and likes/dislikes on posts and comments.

## How it works

`main.go` initializes SQLite, creates tables, seeds categories, and registers HTTP handlers. Go templates render the pages on the server. Registration/sign-in code handles accounts and bcrypt password checks; sessions are associated with a `session_token` cookie. Form handlers change posts, comments, and reactions through the database package, then return or redirect to server-rendered pages. CSS is served from `Frontend/static/`.

## Code map

- `main.go`: route registration and port 8080 startup.
- `backend/database/`: schema, categories, users, sessions, posts, comments, and reactions.
- `backend/pkg/auth/`: signup, signin, and logout.
- `backend/pkg/Pages/`: page rendering.
- `backend/pkg/handlers/`: post/comment/reaction actions.
- `Frontend/templates/`, `Frontend/static/`: HTML and CSS.
- `dockerfile`: multi-stage image build.

## Run locally

Run these commands from this project’s folder.

Requires Go 1.23 or newer and a C compiler with CGO enabled for `github.com/mattn/go-sqlite3`.

```bash
go mod download
go run .
```

Open `http://localhost:8080`. Alternatively:

```bash
docker build -f dockerfile -t reboot-forum .
docker run --rm -p 8080:8080 reboot-forum
```

## Current behavior and limitations

The database path is `./ForoumDatabase.db`; run from this directory so templates and data resolve correctly. Startup creates the schema and seeds categories. Native runs persist changes in that database file. The Docker example has no persistent volume, so container removal discards container-local database changes. The bundled executable is not needed for `go run` and may target a different system.

## Try it

Create an account, sign in, publish a categorized post, add a comment, and change a reaction. Sign out and verify which views/actions remain available.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
