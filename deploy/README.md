# Deploying population-sim to `sim.elliotkirk.com`

The population sim is hosted on the same DigitalOcean droplet as
[`elliotkirk`](https://github.com/eskirk/elliotkirk) and reuses the same SSH
deploy secrets. It has two halves:

- **client** — a static Vite build (`client/dist`) served directly by nginx.
- **server** — a Go WebSocket binary run as a `systemd --user` service on
  `localhost:8080`, reverse-proxied (as `wss://`) by nginx.

```
browser ──https──> nginx (sim.elliotkirk.com)
                     ├── /        -> static files in client/dist
                     └── /ws      -> ws proxy -> localhost:8080 (go service)
```

## Continuous deploy

`.github/workflows/deploy.yml` runs on every push to `main`. It SSHes into the
droplet (using `SSH_HOST`, `SSH_KEY`, `SSH_USERNAME`, `SSH_PASSPHRASE` — the same
secrets elliotkirk uses), pulls `main`, rebuilds the client (with
`VITE_WS_URL=wss://sim.elliotkirk.com/ws` baked in) and the Go binary, then
restarts the user service.

## One-time server setup

Run these once on the droplet as the deploy user (e.g. `steve`).

1. **Clone the repo** into the home directory (the workflow expects
   `~/population-sim`) with pull access (deploy key or token):
   ```sh
   git clone <repo-url> ~/population-sim
   ```

2. **Toolchains** (if not already present from elliotkirk): `nvm` with Node 24,
   and Go ≥ 1.23 at `/usr/local/go`.

3. **Server service** — install the user unit and allow it to run without an
   active login:
   ```sh
   mkdir -p ~/.config/systemd/user
   cp ~/population-sim/deploy/population-sim.service ~/.config/systemd/user/
   loginctl enable-linger "$USER"
   cd ~/population-sim/server && go build -o ../bin/population-sim . && cd -
   systemctl --user daemon-reload
   systemctl --user enable --now population-sim
   ```

4. **nginx site**:
   ```sh
   sudo cp ~/population-sim/deploy/nginx-sim.elliotkirk.com.conf \
       /etc/nginx/sites-available/sim.elliotkirk.com
   sudo ln -s /etc/nginx/sites-available/sim.elliotkirk.com \
       /etc/nginx/sites-enabled/
   sudo nginx -t && sudo systemctl reload nginx
   ```

5. **DNS + TLS** — point an `A` record `sim.elliotkirk.com` at the droplet, then:
   ```sh
   sudo certbot --nginx -d sim.elliotkirk.com
   ```

6. **First client build** (the deploy workflow does this automatically afterward):
   ```sh
   cd ~/population-sim/client
   yarn install --frozen-lockfile
   VITE_WS_URL="wss://sim.elliotkirk.com/ws" yarn build
   ```

## Configuration reference

| Where | Variable | Purpose |
| --- | --- | --- |
| client build | `VITE_WS_URL` | WebSocket endpoint baked into the bundle (`wss://sim.elliotkirk.com/ws`). |
| server service | `ADDR` | Listen address (`localhost:8080`). |
| server service | `ALLOWED_ORIGINS` | Comma-separated `Origin` allowlist (`https://sim.elliotkirk.com`). |

## Useful commands

```sh
systemctl --user status population-sim     # service state
journalctl --user -u population-sim -f     # follow server logs
```
