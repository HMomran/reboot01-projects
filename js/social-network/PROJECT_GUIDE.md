# Social Network

## What it is

A social platform with a Vue frontend and Go/SQLite backend: accounts, profile privacy, followers, feeds, post visibility, comments/reactions, groups, events, messages, uploads, and notifications.

## How it works

Vue Router selects views and frontend stores share auth, chat, and notification state. `apiFetch` sends credentialed requests through `/api`. During development, Vite removes that prefix and proxies to Go on port 8080; it also proxies `/ws` and `/uploads`. The backend runs SQLite migrations at startup, authenticates sessions, and routes requests to feature handlers. Separate WebSocket hubs deliver private chat, group chat, feed changes, presence, and notifications. Docker serves the built frontend through Nginx.

## Code map

- `frontend/src/views/`, `components/`: screens and reusable interface pieces.
- `frontend/src/router/`, `stores/`: navigation and shared frontend state.
- `frontend/src/services/api.js`: HTTP client and auth-expiry handling.
- `frontend/vite.config.js`, `nginx.conf`: development/production proxy setup.
- `backend/main.go`: API routes and WebSocket hubs.
- `backend/pkg/auth/`, `middleware/`: sessions and request authentication.
- `backend/pkg/handlers/`: social features and upload routes.
- `backend/pkg/websocket/`: realtime delivery.
- `backend/pkg/db/migrations/sqlite/`: versioned schema and seed migrations.
- `docker-compose.yml`, `run.sh`: container setup and helper commands.

## Run locally

Run these commands from this project’s folder.

For Docker, run from this folder:

```bash
docker compose up --build -d
```

Open `http://localhost` (frontend port 80, backend port 8080). Stop with `docker compose down`.

For native development, use Go 1.22 or newer with CGO and a C compiler. Check the installed Vite package’s Node engine requirement before choosing a Node runtime. This copy declares Vite 8, while its frontend Dockerfile uses `node:20-alpine`; the Docker build has not been verified against that dependency combination. In one terminal:

```bash
cd backend
go mod download
go run .
```

In another terminal, from the project folder:

```bash
cd frontend
npm install
npm run dev
```

Open the URL printed by Vite. `npm run build` in `frontend/` builds the frontend.

## Current behavior and limitations

Run the backend from `backend/` so relative migration and upload paths resolve. Native storage defaults to `./social.db`; `DB_PATH` overrides it. Docker stores the database and uploads in named volumes. Startup includes seed migrations. `docker compose down` retains these volumes; `./run.sh clean` explicitly deletes them after confirmation. API access uses session cookies. Several other projects use port 8080, so avoid port conflicts. Dependency/runtime compatibility and services were not executed during this documentation pass.

## Try it

Register two accounts, follow a user, switch profile privacy, create a post, and comment/react from the other account. Create a group and event, then exchange private/group messages and check notifications.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
