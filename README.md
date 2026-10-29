# Reboot01 Projects

Source projects organized by language, with a guide for every project. The `js/` collection also contains apps whose servers are written in Go.

Each guide explains purpose, execution flow, important files, setup, example usage, and known gaps found in the copied source. Existing project READMEs are preserved and link to their guide.

Commands in a guide run from that project’s folder unless stated otherwise. These are source-based setup instructions; applications, credentials, and external APIs were not exercised during the documentation pass.

## C projects

| Project | What it does |
| --- | --- |
| [atm-management-system](C/atm-management-system/PROJECT_GUIDE.md) | A terminal bank-account simulator with registration, account management, deposits, withdrawals, ownership transfers, and text-file storage. |

## Go projects

| Project | What it does |
| --- | --- |
| [Ascii-art](go/Ascii-art/PROJECT_GUIDE.md) | A Go command-line program that turns printable ASCII text into eight-row text banners using the bundled `Found.txt` font. |
| [Ascii-art-fs](go/Ascii-art-fs/PROJECT_GUIDE.md) | A Go banner renderer that lets the caller choose standard, shadow, or thinkertoy lettering. |
| [Ascii-art-justify](go/Ascii-art-justify/PROJECT_GUIDE.md) | This folder is a placeholder for an ASCII-art justification variant. No implementation files were present in the copied source. |
| [Ascii-art-output](go/Ascii-art-output/PROJECT_GUIDE.md) | A Go banner renderer that writes ASCII art into a chosen text file. |
| [Ascii-art-web](go/Ascii-art-web/PROJECT_GUIDE.md) | A Go website where a visitor submits text and selects a banner font to generate ASCII art. |
| [Ascii-art-web-dockerize](go/Ascii-art-web-dockerize/PROJECT_GUIDE.md) | A Go website where a visitor submits text and selects a banner font to generate ASCII art. This variant includes a multi-stage Docker build. |
| [forum](go/forum/PROJECT_GUIDE.md) | A Go and SQLite discussion forum with registration, sign-in, categorized posts, comments, and likes/dislikes on posts and comments. |
| [groupie-tracker](go/groupie-tracker/PROJECT_GUIDE.md) | A server-rendered Go music browser that lists artists and shows band details with concert dates grouped by location. |
| [lem-in](go/lem-in/PROJECT_GUIDE.md) | A partially copied Go graph project intended to route ants from a start room to an end room and simulate their movement. |
| [math-skills](go/math-skills/PROJECT_GUIDE.md) | A Go utility that reads one number per line from a text file and prints the average, median, population variance, and standard deviation. |
| [net-cat](go/net-cat/PROJECT_GUIDE.md) | A Go TCP chat server that multiple terminal clients can join using Netcat. |

## JavaScript projects

| Project | What it does |
| --- | --- |
| [MAKE-GAME](js/MAKE-GAME/PROJECT_GUIDE.md) | A browser game built with HTML, CSS, and native JavaScript modules. Move a paddle, bounce a ball into bricks, and complete three timed rounds before running out of lives. |
| [graphql](js/graphql/PROJECT_GUIDE.md) | A browser dashboard for a Reboot01 learner account, showing identity, XP progress, audit ratio, skills, module selection, and recent project activity. |
| [make-your-game](js/make-your-game/PROJECT_GUIDE.md) | A browser game built with HTML, CSS, and native JavaScript modules. Move a paddle, bounce a ball into bricks, and complete three timed rounds before running out of lives. |
| [mini-framework](js/mini-framework/PROJECT_GUIDE.md) | A dependency-free JavaScript UI framework with declarative element trees, DOM rendering, shared state, and hash routing. A TodoMVC app demonstrates its public API. |
| [real-time-forum](js/real-time-forum/PROJECT_GUIDE.md) | A discussion app with a plain JavaScript frontend and Go/SQLite backend. It provides accounts, posts, comments, voting, uploads, private messages, and online-user updates. |
| [real-time-forum-typing-in-progress](js/real-time-forum-typing-in-progress/PROJECT_GUIDE.md) | A discussion app with a plain JavaScript frontend and Go/SQLite backend. It provides accounts, posts, comments, voting, uploads, private messages, and online-user updates. This variant also sends live chat typing indicators. |
| [social-network](js/social-network/PROJECT_GUIDE.md) | A social platform with a Vue frontend and Go/SQLite backend: accounts, profile privacy, followers, feeds, post visibility, comments/reactions, groups, events, messages, uploads, and notifications. |

## Setup notes

- Each Go module is independent; run Go commands in the folder containing its `go.mod`.
- SQLite projects using `mattn/go-sqlite3` require CGO and a C compiler for native builds.
- Several web projects use port 8080; both realtime forum versions use 5500 by default. Start them separately or adjust their configuration.
- Browser-only projects use native JavaScript modules; serve them over HTTP.
- `go/Ascii-art-justify` has no source implementation in this copy.
- `go/lem-in` is missing the simulation package and route-selection implementation referenced by its entry point.
- ASCII-art font variants contain local font path/capitalization mismatches; their guides explain when the network fallback is used.
