package handlers

import (
	"fmt"

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

	destination := GetClosestFreeTile(world, enemyPos.Position, positionComp.Position, positionComp.BoundingBox, enemyPos.BoundingBox)

	generatedPath := false

	//----------- Calculating path --------------------------
	//If path is 0, or our destination changed
	if len(path) == 0 ||
		(destination.X != movementComp.Target.X && destination.Y != movementComp.Target.Y) {
		p, _, found := world.Grid.GetPath(positionComp.Position, destination, positionComp.BoundingBox)
		path = p
		generatedPath = true
		if !found {
			fmt.Printf("Didnt find path %v\n", ev.Index)
			return
		}
	}

	if !generatedPath && len(path) > 0 { // if we haven't generated a path already

		posToMove := path[len(path)-1]

		//if the cell will be occupied or is taken already
		if world.Grid.IsPositionAvailable(posToMove, positionComp.BoundingBox) {
			p, _, found := world.Grid.GetPath(positionComp.Position, destination, positionComp.BoundingBox)
			path = p
			generatedPath = true
			if !found {
				return
			}
		}
	}

	if len(path) == 0 {
		return
	}
	// -----------------------------------------

	// cell, _ := world.Grid.CellAt(path[len(path)-1])
	// if cell.Flag.OccupiedInSteps != -1 && cell.Flag.OccupiedInSteps <= movementComp.MovementSpeed {
	// 	movementComp.Path = movementComp.Path[:0]
	// 	world.ObjectPool.Components["MovementComponent"][ev.Index] = movementComp
	// 	return
	// }

	// if len(path) > 0 {
	// 	nextCell, _ := world.Grid.CellAt(path[len(path)-1])
	// 	nextCell.FlagCell(movementComp.MovementSpeed)
	// }

	//Adding the client event for movement directly
	data := make(map[string]interface{}, 3)
	data["event"] = "walk"
	data["who"] = h.World.EntityManager.GetEntities()[ev.Index].ID
	data["where"] = path[len(path)-1]
	h.World.ClientEventManager.AddEvent(data)

	movementComp.Path = path
	movementComp.IsMoving = true
	movementComp.TimeSinceLastMovement = world.Tick
	world.ObjectPool.Components["MovementComponent"][ev.Index] = movementComp
}

func GetClosestFreeTile(world *game.World, unitPos engine.Vector, myPos engine.Vector, my_bbox engine.Vector, enemybbox engine.Vector) engine.Vector {

	closestFreeTile := engine.Vector{}
	closestDistance := 100000
	tiles := world.Grid.GetSurroundingTilesWithOffset(unitPos, my_bbox.X/2+enemybbox.X/2) // TODO:  nonsquare objects
	for _, tile := range tiles {
		if world.Grid.IsCellTaken(tile) {
			continue
		}
		d := world.Grid.GetDistance(tile, myPos)
		if d < closestDistance {
			closestFreeTile = tile
			closestDistance = d
		}
	}

	return closestFreeTile
}
