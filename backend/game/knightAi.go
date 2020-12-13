package game

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type KnightAI struct {
	world *World
}

func (ai KnightAI) CalculateAction(index int) engine.Action {
	w := ai.world

	atkComp := w.ObjectPool.Components["AttackComponent"][index].(AttackComponent)
	posComp := w.ObjectPool.Components["PositionComponent"][index].(PositionComponent)
	movComp := w.ObjectPool.Components["MovementComponent"][index].(MovementComponent)

	canAttack := false
	canMove := false

	if w.Tick-atkComp.TimeSinceLastAttack >= atkComp.AttackSpeed {
		canAttack = true
	} else {
		//		fmt.Printf("Cannot attack %v, on cooldown for %v\n", index, atkComp.AttackSpeed-(w.Tick-atkComp.TimeSinceLastAttack))
	}

	if w.Tick-movComp.TimeSinceLastMovement >= movComp.MovementSpeed {
		canMove = true
	} else {
		//	fmt.Printf("Cannot move %v, on cooldown for %v\n", index, movComp.MovementSpeed-(w.Tick-movComp.TimeSinceLastMovement))
	}

	if !canAttack && !canMove {
		//	fmt.Printf("Cant move or attack %v\n", index)
		return EmptyAction{}
	}

	nearbyEntities := GetNearbyEntities(40, w, index)

	//If we're already attacking, keep attacking
	if atkComp.Target != -1 && w.EntityManager.Entities[atkComp.Target].Active {
		tarPos := w.ObjectPool.Components["PositionComponent"][atkComp.Target].(PositionComponent)
		if w.Grid.GetDistanceIncludingDiagonal(posComp.Position, tarPos.Position) < 2 {
			if canAttack {
				return AttackAction{Target: atkComp.Target, Index: index}
			}
			return EmptyAction{}
		}
	}

	//Check if an enemy is in range or move to somewhere
	closestIndex := -1
	closestDistance := 100000
	for _, indx := range nearbyEntities {
		if w.EntityManager.Entities[index].PlayerTag != w.EntityManager.Entities[indx].PlayerTag {
			tarPos := w.ObjectPool.Components["PositionComponent"][indx].(PositionComponent)
			if w.Grid.GetDistanceIncludingDiagonal(tarPos.Position, posComp.Position) < 2 {
				if canAttack {
					return AttackAction{Target: indx, Index: index}
				}

				return EmptyAction{}
			}

			dist := w.Grid.GetDistance(tarPos.Position, posComp.Position)
			if dist < closestDistance {
				closestIndex = indx
				closestDistance = dist
			}

		}
	}

	if canAttack {
		//Reset target to noone
		atkComp.Target = -1
		w.ObjectPool.Components["AttackComponent"][index] = atkComp
	}
	if !canMove {
		return EmptyAction{}
	}

	return MovementAction{Target: closestIndex, Index: index}
}
