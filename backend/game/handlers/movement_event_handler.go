package handlers

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

//MovementEventHandler is a handler used to handle Movement of the entities, Handles the MovementAction
//Calculates the neX)t position an entity should be at
type MovementEventHandler struct {
	World *game.World
}

//HandleEvent handles Movement Events for entities, it doesn't move the entity, it just creates the path and marks the neX)t grid as occupied
//
func (h MovementEventHandler) HandleEvent(ev engine.Event) {
	//	world := h.World
	if ev.ID != "MovementEvent" {
		panic(fmt.Sprintf("MovementEventHandler got event other than movement event, %v", ev.Index))
	}

	movementComp := h.World.GetObjectPool().Components["MovementComponent"][ev.Index].(components.MovementComponent)
	positionComp := h.World.ObjectPool.Components["PositionComponent"][ev.Index].(components.PositionComponent)

	movementComp.IsMoving = true

	positionComponents := h.World.ObjectPool.Components["PositionComponent"]

	closestIndex := -1
	closestDistance := float32(100000)

	opposingTag := 0
	if ev.Data["tag"].(int) == 0 {
		opposingTag = 1
	}

	entities := helper.GetTaggedNearbyEntities(10000, h.World, ev.Index, opposingTag)

	for _, entIndex := range entities {
		tarPos := positionComponents[entIndex].(components.PositionComponent)
		dist := engine.GetDistanceIncludingDiagonal(tarPos.Position, positionComp.Position)
		if dist < closestDistance {
			closestIndex = entIndex
			closestDistance = dist
		}
	}

	movementComp.DesiredDirection = positionComponents[closestIndex].(components.PositionComponent).Position.Subtract(positionComp.Position).Normalize()

	data := make(map[string]interface{}, 3)
	data["event"] = "walk"
	data["who"] = h.World.EntityManager.GetEntities()[ev.Index].ID
	data["where"] = positionComp.Position
	data["velocity"] = movementComp.Velocity
	h.World.ClientEventManager.AddEvent(data)

	h.World.GetObjectPool().Components["MovementComponent"][ev.Index] = movementComp
}
