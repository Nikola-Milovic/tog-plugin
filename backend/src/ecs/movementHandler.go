package ecs

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/src/action"
)

//MovementHandler is a handler used to handle Movement of the entities, Handles the MovementAction
//Calculates the next position an entity should be at
type MovementHandler struct {
	manager *EntityManager
}

//HandleAction handles Movement Action for entity at the given index
func (h MovementHandler) HandleAction(index int) {
	action, ok := h.manager.Actions[index].(action.MovementAction)

	if !ok {
		fmt.Println("Error")
	}
	destination := action.Target

	path := h.manager.MovementComponents[index].Path

	//fmt.Printf("Target is %v \n", destination)

	if len(path) == 0 {
		fmt.Println("Calculating path")
		p, dist, found := h.manager.Grid.GetPath(h.manager.PositionComponents[index].Position, destination)
		path = p
		if !found {
			return
		}

		fmt.Printf("Path is %v, and distance %v \n", path, dist)
	}

	h.manager.MovementComponents[index].Path = path

}
