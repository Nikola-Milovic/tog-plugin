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

	generatedPath := false

	if len(path) == 0 ||
		(destination.X != movementComp.Target.X && destination.Y != movementComp.Target.Y) {
		p, _, found := h.world.Grid.GetPath(positionComp.Position, destination)
		path = p
		generatedPath = true
		if !found {
			return
		}
	}

	if !generatedPath && len(path) > 0 {

		posToMove := path[len(path)-1]
		cell, _ := h.world.Grid.CellAt(posToMove)

		if (h.world.Grid.IsCellTaken(posToMove)) ||
			(cell.Flag.OccupiedInSteps != -1 && cell.Flag.OccupiedInSteps <= movementComp.Speed-movementComp.CanMove) {
			p, _, found := h.world.Grid.GetPath(positionComp.Position, destination)
			path = p
			generatedPath = true
			if !found {
				return
			}
		}
	}

	if movementComp.CanMove > 0 {
		fmt.Printf("Can move %v for index %v\n", movementComp.CanMove, action.Index)
		return
	}

	//fmt.Printf("Target is %v \n", destination)

	fmt.Printf("Moving %v \n", action.Index)

	h.world.Grid.ReleaseCell(positionComp.Position)

	posToMove := path[len(path)-1]

	positionComp.Position = posToMove

	h.world.Grid.OccupyCell(posToMove)

	path = path[:len(path)-1]

	if len(path) > 0 {
		nextCell, _ := h.world.Grid.CellAt(path[len(path)-1])
		nextCell.FlagCell(movementComp.Speed)
	}

	movementComp.Path = path
	movementComp.CanMove = movementComp.Speed

	h.world.ObjectPool.Components["MovementComponent"][action.Index] = movementComp
	h.world.ObjectPool.Components["PositionComponent"][action.Index] = positionComp
}
