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
			if int(dist) <= e.AttackComponents[index].Range+10 {
				return action.AttackAction{Target: ind}
			}
			if dist < min {
				target = ind
				min = dist
			}
		}

	}
	fmt.Printf("I at %v , found target at %v at distance %v \n", index, target, e.PositionComponents[target].Position.Distance(e.PositionComponents[index].Position))

	return action.MovementAction{Destination: e.PositionComponents[target].Position}
}
