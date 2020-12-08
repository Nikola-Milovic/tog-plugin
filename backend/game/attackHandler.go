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

	h.world.ObjectPool.Components["AttackComponent"][action.Index] = attackComp

	fmt.Printf("Range is %v\n", attackComp.Range)
}
