package ecs

import (
	"github.com/Nikola-Milovic/tog-plugin/src/action"
)

type AI interface {
	CalculateAction(index int, e *EntityManager) action.Action
}

type KnightAI struct{}

func (ai KnightAI) CalculateAction(index int, e *EntityManager) action.Action {

	nearby := e.getNearbyEntities(20, e.PositionComponents[index].Position, index)

	return action.MovementAction{Target: e.PositionComponents[nearby[0]].Position}
}
