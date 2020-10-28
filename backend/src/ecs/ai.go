package ecs

import (
	"github.com/Nikola-Milovic/tog-plugin/src/action"
)

type AI interface {
	CalculateAction(index int, e *EntityManager) action.Action
}

type KnightAI struct{}

func (ai KnightAI) CalculateAction(index int, e *EntityManager) action.Action {

	target := 0

	min := 10000.0
	for ind := range e.getNearbyEntities(1000, e.PositionComponents[index].Position) {

		if e.Entities[ind] != e.Entities[index] {
			dist := e.PositionComponents[ind].Position.Distance(e.PositionComponents[index].Position)
			if dist < min {
				target = ind
				min = dist
			}
		}
	}

	return action.MovementAction{Destination: e.PositionComponents[target].Position}
}
