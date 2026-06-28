package environment

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// Environment holds the shared world state. All access to actors is guarded by
// mu because the simulation ticker, state reads, and grab handling run on
// separate goroutines.
type Environment struct {
	mu     sync.RWMutex
	actors []Actor
	height int32
	width  int32
}

func SetupEnvironment() *Environment {
	height := int32(1000)
	width := int32(2000)
	size := rand.Int31n(50)

	actors := make([]Actor, 0, size)

	for i := 0; i < int(size); i++ {
		actors = append(actors, Actor{fmt.Sprintf("%d", i), rand.Int31n(width), rand.Int31n(height), false})
	}

	return &Environment{actors: actors, height: height, width: width}
}

func (e *Environment) Run() {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for range ticker.C {
		e.Tick()
	}
}

func (e *Environment) Tick() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i := range e.actors {
		if e.actors[i].grabbed {
			continue
		}

		e.actors[i].move(e)
	}
}

// getActor returns the actor with the given name, or nil if the name is not a
// valid index. Callers must hold the lock.
func (e *Environment) getActor(name string) *Actor {
	ndx, err := strconv.Atoi(name)
	if err != nil || ndx < 0 || ndx >= len(e.actors) {
		return nil
	}

	return &e.actors[ndx]
}

func (e *Environment) GrabActor(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	a := e.getActor(name)
	if a == nil {
		return
	}

	a.grabbed = !a.grabbed
}

// GetState returns the current world as a JSON-encoded array of actors.
func (e *Environment) GetState() []byte {
	e.mu.RLock()
	defer e.mu.RUnlock()

	output := make([]interface{}, 0, len(e.actors))

	for _, actor := range e.actors {
		output = append(output, map[string]interface{}{
			"name":      actor.name,
			"positionX": actor.positionX,
			"positionY": actor.positionY,
			"grabbed":   actor.grabbed,
		})
	}

	out, _ := json.Marshal(output)

	return out
}
