package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//AttackEventHandler is a handler used to handle Attacking, WIP
type AttackEventHandler struct {
	World *World
}

//HandleEvent handles Attack Event for entity at the given index
func (h AttackEventHandler) HandleEvent(ev engine.Event) {

	if ev.ID != constants.AttackEvent {
		panic(fmt.Sprint("Got wrong type of event in AttackEventHandler"))
	}

	attackComp := h.World.ObjectPool.Components["AttackComponent"][ev.Index].(components.AttackComponent)

	target := ev.Data["target"].(int)
	attackComp.Target = target
	attackComp.TimeSinceLastAttack = h.World.Tick

	takeDamageEvent := engine.Event{}
	takeDamageEvent.ID = constants.TakeDamageEvent
	takeDamageEvent.Index = ev.Index
	takeDamageEvent.Priority = constants.AttackEventPriority
	data := make(map[string]interface{}, 3)
	data["index"] = target
	data["amount"] = attackComp.Damage
	data["type"] = "physical"
	takeDamageEvent.Data = data

	//	fmt.Printf("Send attack event, %v is attacking %v\n", ev.Index, target)

	h.World.EventManager.SendEvent(takeDamageEvent)

	fmt.Printf("Attack %v\n", ev.Index)

	if attackComp.OnHit != "" {
		h.onHitEvent(ev, attackComp.OnHit)
	}

	h.World.ObjectPool.Components["AttackComponent"][ev.Index] = attackComp
}

func (h AttackEventHandler) onHitEvent(ev engine.Event, effect string) {
	//attackComp := h.World.ObjectPool.Components["AttackComponent"][ev.Index].(components.AttackComponent)
	event := engine.Event{}
	event.ID = constants.ApplyEffectEvent
	event.Index = ev.Index
	event.Priority = constants.ApplyEffectEventPriority
	data := make(map[string]interface{}, 1)
	data["effectID"] = effect
	data["target"] = ev.Data["target"].(int)
	event.Data = data

	h.World.EventManager.SendEvent(event)
}
