# Brick Breaker — MAKE-GAME

## What it is

A browser game built with HTML, CSS, and native JavaScript modules. Move a paddle, bounce a ball into bricks, and complete three timed rounds before running out of lives.

## How it works

`index.html` provides the game elements and loads `js/index.js`. Input handlers update a `Game` instance. Its `requestAnimationFrame` loop calculates elapsed-frame movement, moves the paddle and ball, checks wall/paddle/brick collisions, and updates DOM positions. Destroyed bricks award points; clearing a round increases speed. Game results are saved in browser local storage and displayed in a top-ten leaderboard ordered by score, then time.

## Code map

- `js/index.js`: startup, keyboard input, and menu buttons.
- `js/game.js`: game state, animation loop, physics, collisions, and round transitions.
- `js/config.js`: dimensions, speeds, three lives, and round time limits.
- `js/ui.js`: menus, score/time displays, notifications, and leaderboard.
- `js/storage.js`: player results under `brickBreakerPlayers` in localStorage.
- `index.html`, `style.css`: DOM and presentation.

## Run locally

Run these commands from this project’s folder.

Requires a modern browser and an HTTP file server; there are no npm dependencies. With Python 3:

```bash
python3 -m http.server 8000
```

Open `http://localhost:8000`. Serve over HTTP so JavaScript modules load correctly.

## Current behavior and limitations

Controls: Left/Right arrows move, Up launches the ball, Escape pauses/resumes. Round limits are 180, 150, and 120 seconds; the game starts with three lives. Results are local to the browser origin, with no online accounts or shared leaderboard. This is a separate version from `../make-your-game`; config, game, and UI files differ, so neither folder should be treated as an exact duplicate. Using the same server origin for both versions shares the same localStorage key. Recorded elapsed time comes from game state, which resets between rounds.

## Try it

Enter a name, launch the ball, break bricks, pause/resume, and restart. Finish or lose a game, then reload to see the locally stored leaderboard.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
