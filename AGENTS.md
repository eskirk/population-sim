# AGENTS.md

## Cursor Cloud specific instructions

`population-sim` has two services that must both run for the app to work end to end:

| Service | Path | Port | Run command |
| --- | --- | --- | --- |
| Go WebSocket server | `server/` | `localhost:8080` | `go run main.go` |
| Vite React + PixiJS client | `client/` | `5173` | `yarn dev` |

### Server (Go)
- Requires **Go >= 1.23** (`go.mod` has `go 1.23.2`). The system Go is 1.22; a 1.23.x
  toolchain is installed at `/usr/local/go` (the startup script puts it on `PATH`). If
  `go version` reports 1.22, prepend `export PATH=/usr/local/go/bin:$PATH`.
- Run from `server/`: `go run main.go`. It streams actor positions as JSON over the
  WebSocket every ~50ms. `CheckOrigin` requires a non-empty `Origin` header, so plain HTTP
  GETs without an `Origin` are rejected — connect with a real WebSocket client/browser.
- Deps: `go mod download`. Vet/build: `go vet ./...`, `go build ./...`.

### Client (Vite)
- Package manager is **yarn**. Both `yarn.lock` and `package-lock.json` are committed, but
  `yarn.lock` is the active one (most recently updated) — use `yarn`, not `npm`.
- Dev: `yarn dev` (port 5173). Build: `yarn build` (`tsc -b && vite build`). Lint: `yarn lint`.
- The WebSocket URL is hard-coded to `127.0.0.1:8080` in `src/App.tsx` (`WS_URL`). The Go
  server must be running first, otherwise the canvas renders empty (no actors).
