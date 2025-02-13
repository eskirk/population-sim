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

	_, message, err := c.ReadMessage()
	newMessage := true

	for {
		go func() {
			_, message, err = c.ReadMessage()

			if err != nil {
				log.Printf("Error %s when reading message from client", err)
			}

			newMessage = true
		}()

		if err != nil {
			log.Printf("Error %s when reading message from client", err)
			return
		}

		if newMessage {
			if strings.HasPrefix(string(message), "grabbed") {
				log.Print(string(message))

				name := strings.Split(string(message), " ")[1]
				wsh.environment.GrabActor(name)
			}

			if strings.HasPrefix(string(message), "mouse") {
				log.Print(string(message))

				x := strings.Split(string(message), " ")[1]
				y := strings.Split(string(message), " ")[2]

				log.Printf("mouse %s %s", x, y)
			}

			newMessage = false
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
				return len(origin) > 0
			},
		},
	}
	http.Handle("/", webSocketHandler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
