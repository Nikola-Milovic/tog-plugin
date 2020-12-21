package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

//AttackEventHandler is a handler used to handle Attacking, WIP
type AttackEventHandler struct {
	world *World
}

//HandleEvent handles Attack Event for entity at the given index
func (h AttackEventHandler) HandleEvent(ev engine.Event) {

	if ev.ID != constants.AttackEvent {
		panic(fmt.Sprint("Got wrong type of event in AttackEventHandler"))
	}

	attackComp := h.world.ObjectPool.Components["AttackComponent"][ev.Index].(AttackComponent)

	target := ev.Data["target"].(int)
	attackComp.Target = target
	attackComp.TimeSinceLastAttack = h.world.Tick

	enemyHealth := h.world.ObjectPool.Components["HealthComponent"][target].(HealthComponent)

	enemyHealth.Health -= attackComp.Damage

	h.world.ObjectPool.Components["AttackComponent"][ev.Index] = attackComp
	h.world.ObjectPool.Components["HealthComponent"][target] = enemyHealth

	//emit take damage event

	//fmt.Printf("Health of %v is %v after attack from %v\n", action.Target, enemyHealth.Health, action.Index)
}
