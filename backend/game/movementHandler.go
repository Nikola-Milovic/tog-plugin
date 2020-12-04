package game

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

//MovementHandler is a handler used to handle Movement of the entities, Handles the MovementAction
//Calculates the next position an entity should be at
type MovementHandler struct {
	world *World
}

//HandleAction handles Movement Action for entity at the given index
func (h MovementHandler) HandleAction(act engine.Action) {
}
