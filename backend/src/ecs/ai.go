package ecs

import (
	"github.com/Nikola-Milovic/tog-plugin/src/action"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

type AI interface {
	CalculateAction(index int, e *EntityManager) action.Action
}

type KnightAI struct{}

func (ai KnightAI) CalculateAction(index int, e *EntityManager) action.Action {

	target := constants.Zero()

	min := float32(10000.0)

	nearby := e.getNearbyEntities(400, e.PositionComponents[index].Position, index)

	for _, ind := range nearby {

		dist := e.PositionComponents[ind].Position.Distance(e.PositionComponents[index].Position)

		if e.Entities[ind].PlayerTag != e.Entities[index].PlayerTag {
			if int(dist) <= e.AttackComponents[index].Range+12 {
				//fmt.Printf("I at %v am attacking %v \n", index, ind)
				return action.AttackAction{Target: ind}
			}

			dest := e.getFreePositionAroundEntity(ind, e.PositionComponents[index].Position)

			if dest != constants.Zero() {
				if dist < min {
					target = dest
				}
			}
		}

	}
	//	fmt.Printf("I at %v , found target at %v at distance %v \n", index, target, e.PositionComponents[target].Position.Distance(e.PositionComponents[index].Position))
	//fmt.Printf("I at %v am walking towards %v \n", index, target)
	return action.MovementAction{Target: target}
}
