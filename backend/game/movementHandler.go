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
	world := h.world
	if !ok {
		panic(fmt.Sprintf("Movement handler got handles action other than movement action, %v", act.GetActionType()))
	}

	movementComp := world.ObjectPool.Components["MovementComponent"][action.Index].(MovementComponent)
	positionComp := world.ObjectPool.Components["PositionComponent"][action.Index].(PositionComponent)
	path := movementComp.Path

	enemyPos := world.ObjectPool.Components["PositionComponent"][action.Target].(PositionComponent)

	destination := getClosestTileToUnit(world, enemyPos.Position, positionComp.Position)

	generatedPath := false

	if len(path) == 0 ||
		(destination.X != movementComp.Target.X && destination.Y != movementComp.Target.Y) {
		p, _, found := world.Grid.GetPath(positionComp.Position, destination)
		path = p
		generatedPath = true
		if !found {
			fmt.Printf("Didnt find path %v\n", action.Index)
			return
		}
	}

	if !generatedPath && len(path) > 0 {

		posToMove := path[len(path)-1]
		cell, _ := world.Grid.CellAt(posToMove)

		if (world.Grid.IsCellTaken(posToMove)) ||
			(cell.Flag.OccupiedInSteps != -1 && cell.Flag.OccupiedInSteps <= movementComp.MovementSpeed-(world.Tick-movementComp.TimeSinceLastMovement)) {
			p, _, found := world.Grid.GetPath(positionComp.Position, destination)
			path = p
			generatedPath = true
			if !found {
				return
			}
		}
	}

	fmt.Printf("Moving %v \n", action.Index)

	world.Grid.ReleaseCell(positionComp.Position)

	posToMove := path[len(path)-1]

	positionComp.Position = posToMove

	world.Grid.OccupyCell(posToMove)

	path = path[:len(path)-1]

	if len(path) > 0 {
		nextCell, _ := world.Grid.CellAt(path[len(path)-1])
		nextCell.FlagCell(movementComp.MovementSpeed)
	}

	movementComp.Path = path
	movementComp.TimeSinceLastMovement = world.Tick

	world.ObjectPool.Components["MovementComponent"][action.Index] = movementComp
	world.ObjectPool.Components["PositionComponent"][action.Index] = positionComp
}

func getClosestTileToUnit(world *World, unitPos engine.Vector, myPos engine.Vector) engine.Vector {

	closestFreeTile := engine.Vector{}
	closestDistance := 100000
	tiles := world.Grid.GetSurroundingTiles(unitPos)
	for _, tile := range tiles {
		d := world.Grid.GetDistance(tile, myPos)
		if d < closestDistance {
			closestFreeTile = tile
			closestDistance = d
		}
	}

	return closestFreeTile
}
