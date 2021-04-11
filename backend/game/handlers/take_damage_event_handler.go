package handlers

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//TakeDamageEventHandler is
type TakeDamageEventHandler struct {
	World *game.World
}

//HandleEvent handles Attack Event for entity at the given index
func (h TakeDamageEventHandler) HandleEvent(ev engine.Event) {

	if ev.ID != constants.TakeDamageEvent {
		panic(fmt.Sprint("Got wrong type of event in TakeDamageEventHandler"))
	}

	target := h.World.EntityManager.GetIndexMap()[ev.Data["target"].(string)]
	amount := ev.Data["amount"].(int)

	health := h.World.ObjectPool.Components["StatsComponent"][target].(components.StatsComponent)

	health.Health -= amount

	//Event for clients
	data := make(map[string]interface{}, 3)
	data["event"] = "take_damage"
	data["who"] = ev.Data["target"]
	data["amount"] = amount
	h.World.ClientEventManager.AddEvent(data)

	fmt.Printf("Take damage %v, amount %v, type %v\n", target, amount, ev.Data["type"].(string))

	h.World.ObjectPool.Components["StatsComponent"][target] = health
}
