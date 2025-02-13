package environment

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Environment struct {
	actors []Actor
	height int32
	width  int32
}

func SetupEnvironment() Environment {
	height := int32(1000)
	width := int32(2000)
	size := rand.Int31n(50)

	actors := make([]Actor, 0, size)

	for i := 0; i < int(size); i++ {
		// log.Print("appending actor")
		actors = append(actors, Actor{fmt.Sprintf("%d", i), rand.Int31n(width), rand.Int31n(height), false})
	}

	return Environment{actors, height, width}
}

func (e Environment) Run() {
	for {
		e.Tick()
		time.Sleep(time.Millisecond * 50)
	}
}

func (e Environment) Tick() {
	// log.Printf("actors %d", len(e.actors))

	for i := range e.actors {
		// log.Print("ticking")

		if e.actors[i].grabbed {
			continue
		}

		e.actors[i].move(e)
	}
}

func (e Environment) getActor(name string) *Actor {
	ndx, err := strconv.Atoi(name)

	if err != nil {
		log.Print("error converting name to int")
	}

	return &e.actors[ndx]
}

func (e Environment) GrabActor(name string) {
	a := e.getActor(name)
	a.grabbed = !a.grabbed

	log.Print("grabbed actor " + a.ToString())
}

func (e Environment) ToString() string {
	state := e.GetState()

	log.Print("actors: " + string(state))

	return string(state)
}

func (e Environment) GetState() []byte {
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
