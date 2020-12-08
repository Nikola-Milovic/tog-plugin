package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
)

//MovementHandler is a handler used to handle Movement of the entities, Handles the MovementAction
//Calculates the next position an entity should be at
type MovementHandler struct {
	world *World
}

//HandleAction handles Movement Action for entity at the given index
func (h MovementHandler) HandleAction(act engine.Action) {
	action, ok := act.(MovementAction)

	if !ok {
		panic(fmt.Sprintf("Movement handler got handles action other than movement action, %v", act.GetActionType()))
	}

	destination := action.Target

	movementComp := h.world.ObjectPool.Components["MovementComponent"][action.Index].(MovementComponent)
	positionComp := h.world.ObjectPool.Components["PositionComponent"][action.Index].(PositionComponent)

	path := movementComp.Path
	//fmt.Printf("Target is %v \n", destination)

	if len(path) == 0 {
		fmt.Println("Calculating path")
		p, dist, found := h.world.Grid.GetPath(positionComp.Position, destination)
		path = p
		if !found {
			return
		}

		fmt.Printf("Path is %v, and distance %v \n", path, dist)
	}

	if h.world.Counter%movementComp.Speed == 0 && len(path) > 0 {
		fmt.Printf("Move %v\n", action.Index)
		positionComp.Position = path[len(path)-1]
		path = path[:len(path)-1]
		fmt.Printf("Position after move is %v\n", positionComp.Position)
	}

	movementComp.Path = path
	h.world.ObjectPool.Components["MovementComponent"][action.Index] = movementComp
	h.world.ObjectPool.Components["PositionComponent"][action.Index] = positionComp

}
