package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//TakeDamageEventHandler is
type TakeDamageEventHandler struct {
	World *World
}

//HandleEvent handles Attack Event for entity at the given index
func (h TakeDamageEventHandler) HandleEvent(ev engine.Event) {

	if ev.ID != constants.TakeDamageEvent {
		panic(fmt.Sprint("Got wrong type of event in TakeDamageEventHandler"))
	}

	target := ev.Data["index"].(int)
	amount := ev.Data["amount"].(int)

	health := h.World.ObjectPool.Components["StatsComponent"][target].(components.StatsComponent)

	health.Health -= amount

	fmt.Printf("Take damage %v, amount %v\n", target, amount)

	h.World.ObjectPool.Components["StatsComponent"][target] = health
}
