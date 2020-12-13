package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
)

//AttackHandler is a handler used to handle Attacking, WIP
type AttackHandler struct {
	world *World
}

//HandleAction handles Attack Action for entity at the given index
func (h AttackHandler) HandleAction(act engine.Action) {
	action, ok := act.(AttackAction)
	if !ok {
		panic(fmt.Sprint("Got wrong type of action in AttackHandler"))
	}

	attackComp := h.world.ObjectPool.Components["AttackComponent"][action.Index].(AttackComponent)

	attackComp.Target = action.Target
	attackComp.TimeSinceLastAttack = h.world.Tick

	enemyHealth := h.world.ObjectPool.Components["HealthComponent"][action.Target].(HealthComponent)

	enemyHealth.Health -= attackComp.Damage

	h.world.ObjectPool.Components["AttackComponent"][action.Index] = attackComp
	h.world.ObjectPool.Components["HealthComponent"][action.Target] = enemyHealth

	fmt.Printf("Health of %v is %v after attack from %v\n", action.Target, enemyHealth.Health, action.Index)
}
