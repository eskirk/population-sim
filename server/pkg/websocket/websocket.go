package websocket

import (
	"log"
	"net/http"
	"strings"
	"time"

	"population-sim/pkg/environment"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader    websocket.Upgrader
	environment environment.Environment
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

	for {
		mt, message, err := c.ReadMessage()
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
		log.Printf("Receive message %s", string(message))
		if strings.Trim(string(message), "\n") != "start" {
			err = c.WriteMessage(websocket.TextMessage, []byte("You did not say the magic word!"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			continue
		}
		log.Println("start responding to client...")
		i := 1

		for {
			wsh.environment.Tick()
			// log.Print(wsh.environment.ToString())

			response := wsh.environment.ToString()
			err = c.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			i = i + 1
			time.Sleep(time.Millisecond * 50)
		}
	}

}

func Serve(env environment.Environment) {
	webSocketHandler := webSocketHandler{
		environment: env,
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
