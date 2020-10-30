package ecs

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/src/action"
)

type AI interface {
	CalculateAction(index int, e *EntityManager) action.Action
}

type KnightAI struct{}

func (ai KnightAI) CalculateAction(index int, e *EntityManager) action.Action {

	target := index

	min := float32(10000.0)

	nearby := e.getNearbyEntities(400, e.PositionComponents[index].Position, index)

	for _, ind := range nearby {

		dist := e.PositionComponents[ind].Position.Distance(e.PositionComponents[index].Position)

		if e.Entities[ind].PlayerTag != e.Entities[index].PlayerTag {
			if int(dist) <= e.AttackComponents[index].Range {
				return action.AttackAction{Target: ind}
			}
			if dist < min {
				target = ind
				min = dist
			}
		}

	}

	return action.MovementAction{Destination: e.PositionComponents[target].Position, NearbyEntities: nearby}
}
