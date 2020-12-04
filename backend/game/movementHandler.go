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
		panic(fmt.Sprintf("Movement handler got handles action other than movement action, %v", act.GetActionState()))
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
	
	if(h.world.Counter % 5 == 0) {
		positionComp.Position = path[len(path)-1]
		path = path[:len(path)-1]
	}

	movementComp.Path = path
	h.world.ObjectPool.Components["MovementComponent"][action.Index] = movementComp
	h.world.ObjectPool.Components["PositionComponent"][action.Index] = positionComp
}
