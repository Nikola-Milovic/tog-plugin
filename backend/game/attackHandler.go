package game

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

//AttackHandler is a handler used to handle Attacking, WIP
type AttackHandler struct {
	manager *engine.EntityManager
}

//HandleAction handles Attack Action for entity at the given index
func (h AttackHandler) HandleAction(act engine.Action) {
}
