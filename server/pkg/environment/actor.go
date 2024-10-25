package environment

import (
	"math/rand"
)

type Actor struct {
	name      string
	positionX int32
	positionY int32
}

func (a Actor) move(e Environment) Actor {
	a.positionX = a.moveDirection(e, "x")
	a.positionY = a.moveDirection(e, "y")

	return a
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
		if a.positionX+movement > environment.height || a.positionX+movement < 0 {
			return a.positionX - movement
		}

		return a.positionX + movement
	}

	if dimension == "y" {
		if a.positionY+movement > environment.width || a.positionY+movement < 0 {
			return a.positionY - movement
		}

		return a.positionY + movement
	}

	return 0
}
