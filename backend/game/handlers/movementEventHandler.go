package handlers

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//MovementEventHandler is a handler used to handle Movement of the entities, Handles the MovementAction
//Calculates the next position an entity should be at
type MovementEventHandler struct {
	World *game.World
}

//HandleEvent handles Movement Events for entities, it doesn't move the entity, it just creates the path and marks the next grid as occupied
//
func (h MovementEventHandler) HandleEvent(ev engine.Event) {
	world := h.World
	if ev.ID != "MovementEvent" {
		panic(fmt.Sprintf("MovementEventHandler got event other than movement event, %v", ev.Index))
	}

	movementComp := world.ObjectPool.Components["MovementComponent"][ev.Index].(components.MovementComponent)
	positionComp := world.ObjectPool.Components["PositionComponent"][ev.Index].(components.PositionComponent)
	path := movementComp.Path

	target := ev.Data["target"].(int)
	enemyPos := world.ObjectPool.Components["PositionComponent"][target].(components.PositionComponent)

	destination := getClosestTileToUnit(world, enemyPos.Position, positionComp.Position)

	generatedPath := false

	//----------- Calculating path --------------------------
	//If something is wrong, recalculate path
	if len(path) == 0 ||
		(destination.X != movementComp.Target.X && destination.Y != movementComp.Target.Y) {
		p, _, found := world.Grid.GetPath(positionComp.Position, destination)
		path = p
		generatedPath = true
		if !found {
			fmt.Printf("Didnt find path %v\n", ev.Index)
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
	// -----------------------------------------

	if len(path) > 0 {
		nextCell, _ := world.Grid.CellAt(path[len(path)-1])
		nextCell.FlagCell(movementComp.MovementSpeed)
	}

	//Event for the client to know where the unit is moving to 
	data := make(map[string]interface{}, 2)
	data["where"] = path[len(path)-1]
	data["emitter"] = h.World.EntityManager.Entities[ev.Index].ID
	clientEv := engine.Event{Index: ev.Index, ID: constants.ClientMovementEvent, Priority: 99, Data: data}
	h.World.ClientEventManager.OnEvent(clientEv)


	movementComp.Path = path
	movementComp.IsMoving = true
	movementComp.TimeSinceLastMovement = world.Tick
	world.ObjectPool.Components["MovementComponent"][ev.Index] = movementComp
}

func getClosestTileToUnit(world *game.World, unitPos engine.Vector, myPos engine.Vector) engine.Vector {

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
