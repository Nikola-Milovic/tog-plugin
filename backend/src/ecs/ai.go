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

	min := 10000.0
	for _, ind := range e.getNearbyEntities(400, e.PositionComponents[index].Position, index) {

		dist := e.PositionComponents[ind].Position.Distance(e.PositionComponents[index].Position)

		if int(dist) <= e.AttackComponents[index].Range {
			return action.AttackAction{Target: ind}
		}
		fmt.Printf("I at index %v , found enemy at %v, at distance %v \n", ind, index, dist)

		if dist < min {
			target = ind
			min = dist
		}

	}

	return action.MovementAction{Destination: e.PositionComponents[target].Position}
}
