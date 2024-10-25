package environment

import (
	"fmt"
	// "log"
	"math/rand"
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
		actors = append(actors, Actor{fmt.Sprintf("%d", i), rand.Int31n(width), rand.Int31n(height)})
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

		e.actors[i] = e.actors[i].move(e)
	}
}

func (e Environment) ToString() string {
	output := ""

	for _, actor := range e.actors {
		// log.Printf("Actor %s x: %d y %d \n", actor.name, actor.positionX, actor.positionY)

		output = output + fmt.Sprintf("%s %d %d\n", actor.name, actor.positionX, actor.positionY)
	}

	return output
}
