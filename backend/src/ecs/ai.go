package ecs

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/src/action"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

type AI interface {
	CalculateAction(index int, e *EntityManager) action.Action
}

type KnightAI struct{}

func (ai KnightAI) CalculateAction(index int, e *EntityManager) action.Action {

	nearby := e.getNearbyEntities(20, e.PositionComponents[index].Position, index)
	//	closestEntity := -1
	closestPosition := 100
	targetPosition := constants.V2{}

	if e.AttackComponents[index].Target != -1 {
		//check if the target we we're attacking is still in range
		if e.Grid.GetDistanceIncludingDiagonal(e.PositionComponents[e.AttackComponents[index].Target].Position, e.PositionComponents[index].Position) <= e.AttackComponents[index].Range {
			fmt.Printf("Target %v is still in range of %v\n", e.AttackComponents[index].Target, index)
			return action.AttackAction{Target: e.AttackComponents[index].Target, Index: index}
		}
	}

	for _, indx := range nearby {
		if e.Entities[indx].PlayerTag == e.Entities[index].PlayerTag {
			continue
		}
		dist := e.Grid.GetDistanceIncludingDiagonal(e.PositionComponents[indx].Position, e.PositionComponents[index].Position)

		if dist <= e.AttackComponents[index].Range {
			//	fmt.Printf("Index %v, is attacking %v\n", index, indx)
			e.AttackComponents[index].Target = indx
			return action.AttackAction{Target: indx, Index: index}
		}

		if dist < closestPosition {
			freePlace := e.getSurroundingFreeCell(0, e.PositionComponents[indx].Position)
			for _, place := range freePlace {
				d := e.Grid.GetDistance(place, e.PositionComponents[index].Position)
				if d < closestPosition {
					//fmt.Println(freePlace)
					targetPosition = place
					//closestEntity = indx
					closestPosition = d
				}
			}

		}
	}
	fmt.Printf("Walking %v\n", index)
	//fmt.Printf("Index %v, is walking towards %v, from position %v to %v\n", index, closestEntity, e.PositionComponents[index].Position, targetPosition)
	return action.MovementAction{Target: targetPosition, Index: index}
}
