package environment

import (
	"fmt"
	// "log"
	"math/rand"
	"time"
)

type Actor struct {
	name      string
	positionX int32
	positionY int32
}

type Environment struct {
	actors []Actor
	height int32
	width  int32
}

func SetupEnvironment() Environment {
	height := int32(800)
	width := int32(1000)
	size := rand.Int31n(50)

	actors := make([]Actor, 0, size)

	for i := 0; i < int(size); i++ {
		// log.Print("appending actor")
		actors = append(actors, Actor{fmt.Sprintf("%d", i), rand.Int31n(height), rand.Int31n(width)})
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

		e.actors[i].positionX = e.actors[i].moveDirection(e, "x")
		e.actors[i].positionY = e.actors[i].moveDirection(e, "y")
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

func (a Actor) moveDirection(environment Environment, dimension string) int32 {
	distance := rand.Int31n(5)
	var direction int32

	if rand.Float32() > 0.5 {
		direction = 1
	} else {
		direction = -1
	}

	movement := distance * direction

	// log.Printf("actor %s dimension %s movement %d env %d %d", a.name, dimension, movement, environment.width, environment.height)

	if dimension == "x" {
		if a.positionX+movement > environment.width {
			return a.positionX - movement
		}

		return a.positionX + movement
	}

	if dimension == "y" {
		if a.positionX+movement > environment.height {
			return a.positionY - movement
		}

		return a.positionY + movement
	}

	return a.positionX + movement
}
