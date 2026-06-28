package websocket

import (
	"log"
	"net/http"
	"os"
	"population-sim/pkg/environment"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// writeWait is how long a single write is allowed to take.
	writeWait = 10 * time.Second
	// pongWait is how long we wait for a pong before considering the peer dead.
	pongWait = 60 * time.Second
	// pingPeriod must be less than pongWait so pings keep the connection alive.
	pingPeriod = (pongWait * 9) / 10
	// statePeriod is how often the full world state is pushed to clients.
	statePeriod = 50 * time.Millisecond
)

type webSocketHandler struct {
	upgrader    websocket.Upgrader
	environment *environment.Environment
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	defer c.Close()

	// done is closed by the single reader goroutine when the connection ends,
	// which signals the writer loop below to stop.
	done := make(chan struct{})

	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// single reader: gorilla/websocket does not allow concurrent reads, so all
	// incoming messages flow through this one goroutine.
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Printf("read error: %v", err)
				}
				return
			}
			wsh.handleMessage(message)
		}
	}()

	stateTicker := time.NewTicker(statePeriod)
	pingTicker := time.NewTicker(pingPeriod)
	defer stateTicker.Stop()
	defer pingTicker.Stop()

	// single writer: only this loop ever writes to the connection.
	for {
		select {
		case <-done:
			return
		case <-stateTicker.C:
			c.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.WriteMessage(websocket.TextMessage, wsh.environment.GetState()); err != nil {
				log.Printf("write error: %v", err)
				return
			}
		case <-pingTicker.C:
			c.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (wsh webSocketHandler) handleMessage(message []byte) {
	msg := string(message)

	switch {
	case strings.HasPrefix(msg, "grabbed"):
		parts := strings.Fields(msg)
		if len(parts) >= 2 {
			log.Print(msg)
			wsh.environment.GrabActor(parts[1])
		}
	case strings.HasPrefix(msg, "mouse"):
		// reserved for future mouse-driven interaction
	}
}

func Serve(env *environment.Environment) {
	allowedOrigins := parseAllowedOrigins()

	webSocketHandler := webSocketHandler{
		environment: env,
		upgrader: websocket.Upgrader{
			CheckOrigin: originChecker(allowedOrigins),
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/", webSocketHandler)

	addr := listenAddr()
	log.Printf("Starting server on %s (allowed origins: %v)", addr, allowedOrigins)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// originChecker validates the request Origin against an allowlist. A "*" entry
// allows any origin (useful for local development behind no proxy).
func originChecker(allowed []string) func(*http.Request) bool {
	return func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return false
		}

		for _, a := range allowed {
			if a == "*" || a == origin {
				return true
			}
		}

		log.Printf("rejected connection from origin: %s", origin)
		return false
	}
}

// parseAllowedOrigins reads a comma-separated ALLOWED_ORIGINS env var, falling
// back to the local Vite dev server origins.
func parseAllowedOrigins() []string {
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		parts := strings.Split(v, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				out = append(out, trimmed)
			}
		}
		return out
	}

	return []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	}
}

// listenAddr reads the ADDR env var, defaulting to localhost:8080.
func listenAddr() string {
	if v := os.Getenv("ADDR"); v != "" {
		return v
	}
	return "localhost:8080"
}
