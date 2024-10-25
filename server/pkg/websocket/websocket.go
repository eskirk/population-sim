package websocket

import (
	"log"
	"net/http"
	"population-sim/pkg/environment"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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
	defer func() {
		log.Printf("closing connection")
		c.Close()
	}()

	mt, message, err := c.ReadMessage()
	latestMessage := string(message)

	for {
		go func() {
			mt, message, err = c.ReadMessage()
		}()

		if err != nil {
			log.Printf("Error %s when reading message from client", err)
			return
		}
		if mt == websocket.BinaryMessage {
			err = c.WriteMessage(websocket.TextMessage, []byte("server doesn't support binary messages"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
			}
			return
		}

		if strings.HasPrefix(string(message), "clicked") && string(message) != string(latestMessage) {
			log.Print(string(message))
			latestMessage = string(message)
		}

		response := wsh.environment.ToString()
		err = c.WriteMessage(websocket.TextMessage, []byte(response))
		if err != nil {
			log.Printf("Error %s when sending message to client", err)
			return
		}
		time.Sleep(time.Millisecond * 50)
	}

}

func Serve(env environment.Environment) {
	webSocketHandler := webSocketHandler{
		environment: &env,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				log.Printf("Origin: %s", origin)
				return origin == "http://127.0.0.1:8080" || origin == "http://localhost:5173"
			},
		},
	}
	http.Handle("/", webSocketHandler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
