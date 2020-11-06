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
func (h MovementHandler) HandleAction(act action.Action) {
	action, ok := act.(action.MovementAction)

	index := action.Index

	if !ok {
		fmt.Println("Error")
	}

	path := h.manager.MovementComponents[index].Path

	destination := action.Target

	if len(path) == 0 || h.manager.Grid.IsCellTaken(path[0].X, path[0].Y) {
		//	fmt.Printf("Calculating path %v \n", index)
		p, _, found := h.manager.Grid.GetPath(h.manager.PositionComponents[index].Position, destination)
		path = p
		if !found {
			return
		}

		//		fmt.Printf("Path is %v, and distance %v \n", path, dist)

	}

	move := false
	h.manager.MovementComponents[index].Tick = h.manager.MovementComponents[index].Tick + 1
	if h.manager.MovementComponents[index].Tick == 5 {
		h.manager.MovementComponents[index].Tick = 0
		move = true
	}

	if move {
		h.manager.Grid.ReleaseCell(h.manager.PositionComponents[index].Position)
		h.manager.PositionComponents[index].Position = path[len(path)-1]
		path = path[:len(path)-1]
		h.manager.Grid.OccupyCell(h.manager.PositionComponents[index].Position)
	}

	h.manager.MovementComponents[index].Path = path

}
