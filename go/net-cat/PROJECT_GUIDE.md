# Net-Cat TCP Chat

## What it is

A Go TCP chat server that multiple terminal clients can join using Netcat.

## How it works

The server accepts TCP connections and starts a goroutine per connection. Clients receive the welcome text, choose a unique printable name, and receive the in-memory message history. Messages are filtered to printable ASCII, timestamped, stored in history, and broadcast to other clients. Mutexes protect the connected-client map and history. Join and leave events are also broadcast.

## Code map

- `main.go`: optional port validation and server startup.
- `function/server.go`: listener, clients, and broadcasting.
- `function/client_handler.go`: names, prompts, and messages.
- `function/History.go`: shared in-memory history.
- `TCP-chat.txt`: welcome banner.

## Run locally

Run these commands from this project’s folder.

Requires Go 1.18 or newer; clients need `nc`.

```bash
go run .
```

In two other terminals:

```bash
nc localhost 8989
```

To choose a different port: `go run . 9000`.

## Current behavior and limitations

The default port is 8989; custom ports must be 1024–65535. Up to 10 named clients are allowed. Names are at most 20 characters after filtering. History is lost on restart. This is raw TCP without account authentication; it binds all interfaces. Reads use a fixed buffer rather than a complete line protocol.

## Try it

Join with two different names, exchange messages, reconnect to see history, and try a duplicate name.

## Documentation scope

This guide describes the source files in this copy. Setup commands and behavior were reviewed against the code; applications and external services were not run as part of this documentation pass.
