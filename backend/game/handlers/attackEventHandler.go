package handlers

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//AttackEventHandler is a handler used to handle Attacking, WIP
type AttackEventHandler struct {
	World *game.World
}

//HandleEvent handles Attack Event for entity at the given index
func (h AttackEventHandler) HandleEvent(ev engine.Event) {

	if ev.ID != constants.AttackEvent {
		panic(fmt.Sprint("Got wrong type of event in AttackEventHandler"))
	}

	attackComp := h.World.ObjectPool.Components["AttackComponent"][ev.Index].(components.AttackComponent)

	target := ev.Data["target"].(int)
	attackComp.Target = h.World.EntityManager.Entities[target].ID
	attackComp.TimeSinceLastAttack = h.World.Tick
	attackComp.IsAttacking = true

	h.World.ObjectPool.Components["AttackComponent"][ev.Index] = attackComp
}
