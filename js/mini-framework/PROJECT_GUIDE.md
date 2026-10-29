# Mini Framework

## What it is

A dependency-free JavaScript UI framework with declarative element trees, DOM rendering, shared state, and hash routing. A TodoMVC app demonstrates its public API.

## How it works

`createElement` creates plain objects describing desired DOM nodes. `render` recursively builds real nodes, applies attributes, attaches event handlers, and replaces the container’s contents. `createStore` holds shared state and notifies subscribers after updates. `createRouter` maps the URL hash to routes and notifies subscribers on navigation. The Todo app combines those pieces: state changes persist todos and redraw the view, while routes select All, Active, or Completed filters.

## Code map

- `framework/createElement.js`: element descriptions and child normalization.
- `framework/render.js`: recursive DOM creation, attributes, and events.
- `framework/store.js`: state access, updates, and subscriptions.
- `framework/router.js`: hash routing and navigation subscriptions.
- `framework/index.js`: public exports.
- `todo/app.js`: todo operations, persistence, routing, and redraws.
- `server.js`: static development server.
- `tests/framework.test.js`: dependency-free automated checks.
- `README.md`: API examples and detailed manual verification steps.

## Run locally

Run these commands from this project’s folder.

Requires Node.js; use Node.js 22 or newer for the included checks. No package installation is needed.

```bash
npm start
```

Open `http://localhost:3000`. Set `PORT` to change the development port.

```bash
npm test
```

## Current behavior and limitations

Rendering replaces children instead of diffing a virtual tree. Applications must preserve or restore input focus/drafts when appropriate, as the Todo example does. Treat store values as immutable. The app uses browser localStorage for persistence and URL hashes for filters. Initial rendering is explicit; subscriptions notify on changes. Automated renderer checks use a DOM substitute, so browser interaction still requires manual verification.

## Try it

Add todos, toggle completion, switch filters, use browser Back/Forward, double-click to edit, clear completed items, and reload. The existing README provides detailed API usage and acceptance checks.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
