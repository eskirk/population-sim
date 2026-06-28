# AGENTS.md

## Cursor Cloud specific instructions

`population-sim` has two services that must both run for the app to work end to end:

| Service | Path | Port | Run command |
| --- | --- | --- | --- |
| Go WebSocket server | `server/` | `localhost:8080` | `go run main.go` |
| Vite React + Three.js client | `client/` | `5173` | `yarn dev` |

### Server (Go)
- Requires **Go >= 1.23** (`go.mod` has `go 1.23.2`). The system Go is 1.22; a 1.23.x
  toolchain is installed at `/usr/local/go` (the startup script puts it on `PATH`). If
  `go version` reports 1.22, prepend `export PATH=/usr/local/go/bin:$PATH`.
- Run from `server/`: `go run main.go`. It streams actor positions as JSON over the
  WebSocket every ~50ms (a single reader goroutine + a ticker-driven writer, with ping/pong
  keepalive). World state is guarded by a mutex.
- Env vars: `ADDR` (listen address, default `localhost:8080`) and `ALLOWED_ORIGINS`
  (comma-separated `Origin` allowlist, default the local Vite dev origins; use `*` to allow
  all). Connections without an allowed `Origin` are rejected, so connect with a real
  WebSocket client/browser sending one.
- Deps: `go mod download`. Vet/build: `go vet ./...`, `go build ./...`. Run the race detector
  with `go run -race main.go` (note: the `-race` build takes ~15s to compile before it binds).

### Client (Vite)
- Package manager is **yarn**. Both `yarn.lock` and `package-lock.json` are committed, but
  `yarn.lock` is the active one (most recently updated) — use `yarn`, not `npm`.
- Dev: `yarn dev` (port 5173). Build: `yarn build` (`tsc -b && vite build`). Lint: `yarn lint`.
- Rendering uses **Three.js** via `@react-three/fiber` + `@react-three/drei` (React 18, so
  fiber is pinned to v8). After changing deps, restart `yarn dev` (clear `node_modules/.vite`)
  so Vite re-optimizes.
- The WebSocket URL comes from `VITE_WS_URL` (see `.env.example`), defaulting to
  `ws://localhost:8080`. The Go server must be running first, otherwise the scene renders the
  empty ground with no actors.
