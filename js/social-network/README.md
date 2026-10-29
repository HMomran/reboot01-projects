# Social Network Application

A full-stack, containerized social network application featuring a Go-based backend and a Vue 3 frontend. The platform supports user profiles, activity feeds, real-time messaging via WebSockets, group creation with event planning, and notification management.

## 🚀 Tech Stack

- **Frontend**: [Vue 3](https://vuejs.org/) (Composition API, `<script setup>`), [Vite](https://vite.dev/), [Vue Router](https://router.vuejs.org/) for routing, and Pinia for state management.
- **Backend**: [Go (Golang) 1.22](https://go.dev/) using standard library HTTP routing, custom middleware (CORS, Auth, Logger), and [Gorilla WebSockets](https://github.com/gorilla/websocket) for real-time features.
- **Database**: [SQLite](https://www.sqlite.org/) with automated schema migrations managed via [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate).
- **Containerization**: Docker & Docker Compose.

---

## 🛠 Features

- **Authentication & Profiles**:
  - Secure registration and login.
  - User profiles with privacy/visibility settings (public/private).
  - Profile customization including avatar details.
- **Social Interaction**:
  - Main activity feed with posts, comments, and reactions.
  - User-to-user followers/following relationships (with request/accept flow for private profiles).
- **Groups & Communities**:
  - Create and join groups.
  - Group-specific posts and discussion boards.
  - Event planning within groups (RSVPs, dates, details).
- **Real-Time Communication**:
  - WebSocket-powered private messaging/chat system.
  - Real-time notification delivery for messages, follows, and invites.
- **Media Management**:
  - Localized file/image uploads for posts and profile pictures.

---

## 📁 Repository Structure

```text
├── backend/                   # Go API codebase
│   ├── pkg/                   # Domain logic packages
│   │   ├── auth/              # JWT tokens, password hashing
│   │   ├── db/                # SQLite connection and migration scripts
│   │   ├── handlers/          # API route handlers (posts, profiles, etc.)
│   │   ├── middleware/        # CORS and authentication filters
│   │   └── utils/             # File upload and helper functions
│   │   └── websocket/         # WebSocket hub and connections
│   ├── main.go                # Backend entrypoint & router initialization
│   ├── Dockerfile             # Container image configuration for backend
│   └── go.mod                 # Go module definition
│
├── frontend/                  # Vue 3 Single Page Application
│   ├── src/
│   │   ├── components/        # Reusable UI components (Navbar, Cards)
│   │   ├── views/             # Full-page views (Feed, Chat, Profile)
│   │   ├── services/          # API fetch and network helpers
│   │   ├── stores/            # Client-side state stores
│   │   └── router/            # Vue Router config
│   ├── Dockerfile             # Nginx container for hosting frontend assets
│   └── package.json           # Node.js dependencies and scripts
│
├── docker-compose.yml         # Local orchestration of services
└── run.sh                     # Helper script for container orchestration
```

---

## ⚙️ Setup and Deployment

To run this application locally, you must use the orchestrator script. Please refer to the dedicated execution guide for instructions:

👉 **[Docker Compose Run Guide (RUN_README.md)](RUN_README.md)**

---

## 👥 Authors

********
********
**Hasan Omran**

---

This project was created as part of **Reboot01**'s full-stack development course.

## Project documentation

[Project guide: purpose, architecture, setup, and usage](PROJECT_GUIDE.md)
