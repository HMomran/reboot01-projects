# Real-Time Forum — Typing Indicators

## What it is

A discussion app with a plain JavaScript frontend and Go/SQLite backend. It provides accounts, posts, comments, voting, uploads, private messages, and online-user updates. This variant also sends live chat typing indicators.

## How it works

The Go server creates/updates `forum.db`, registers JSON endpoints under `/api/`, serves uploaded files, and serves the sibling `frontend/` directory. Browser scripts build the interface, submit HTTP requests for stored data, and open `/ws` for live events. The `db` package handles persistence; handlers validate/authenticate requests and call it. Sessions use backend session records with cookie/token handling, and frontend state is kept in sessionStorage. Chat input emits `typing` and `stop_typing`; the WebSocket handler forwards those transient events to the target user. A 1.5-second idle timer stops the indicator.

## Code map

- `backend/main.go`: HTTP routes, static serving, CORS, and startup.
- `backend/db/`: SQLite tables and user/post/comment/message/session/vote queries.
- `backend/handlers/`: authentication, posts, comments, votes, uploads, and WebSocket handling.
- `backend/models/models.go`: shared data structures.
- `frontend/app.js`: interface construction and common utilities.
- `frontend/Login.js`, `Registration.js`, `Posts.js`, `PostDetail.js`: page behavior.
- `frontend/Chat.js`: chat state and WebSocket events.
- `frontend/init.js`: initial data loading and backend health request.

## Run locally

Run these commands from this project’s folder.

Requires Go 1.25.5 or newer and a C compiler with CGO enabled for SQLite. From this project folder:

```bash
cd backend
go mod download
go run .
```

Open `http://localhost:5500`. The Go server serves both frontend and backend; no frontend build is required.

## Current behavior and limitations

This server listens on fixed port 5500. The frontend uses relative `/api/` requests and opens WebSockets on the current page host, so it should be served by this Go server or by a proxy that forwards those paths. Run the two forum variants separately because both default to port 5500. Database and uploads are relative to `backend/`. This guide describes this copied version without asserting that any deployed service is available.

## Try it

Use two browser profiles with two accounts. Publish a post and comment, exchange private messages, and check online presence. Type in one profile without sending and verify that the other sees the typing indicator clear after input stops.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
