package handlers

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/game/components"

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
	movementComp.IsMoving = true

	positionComp := h.World.ObjectPool.Components["PositionComponent"][ev.Index].(components.PositionComponent)

	data := make(map[string]interface{}, 3)
	data["event"] = "walk"
	data["who"] = h.World.EntityManager.GetEntities()[ev.Index].ID
	data["where"] = positionComp.Position
	data["velocity"] = movementComp.Direction
	h.World.ClientEventManager.AddEvent(data)

	h.World.GetObjectPool().Components["MovementComponent"][ev.Index] = movementComp

}