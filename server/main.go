package main

import (
	"population-sim/pkg/environment"
	"population-sim/pkg/websocket"
)

func main() {
	env := environment.SetupEnvironment()
	go env.Run()
	websocket.Serve(env)
}
