package game

import (
	"fmt"

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

	nearbyEntities := GetNearbyEntities(30, w, index)

	//If we're already attacking, keep attacking
	if atkComp.Target != -1 {
		tarPos := w.ObjectPool.Components["PositionComponent"][atkComp.Target].(PositionComponent)
		if w.Grid.GetDistanceIncludingDiagonal(posComp.Position, tarPos.Position) < 2 {
			if canAttack {
				return AttackAction{Target: atkComp.Target, Index: index}
			}
			fmt.Printf("Still in range %v but cannot attack %v", atkComp.Target, index)
			return EmptyAction{}
		}
	}

	//Check if an enemy is in range or move to somewhere
	closestFreeTile := engine.Vector{}
	closestDistance := 100000
	for _, indx := range nearbyEntities {
		if w.EntityManager.Entities[index].PlayerTag != w.EntityManager.Entities[indx].PlayerTag {
			tarPos := w.ObjectPool.Components["PositionComponent"][indx].(PositionComponent)
			if canAttack && w.Grid.GetDistanceIncludingDiagonal(tarPos.Position, posComp.Position) < 2 {
				return AttackAction{Target: atkComp.Target, Index: index}
			}

			if canMove {
				tiles := w.Grid.GetSurroundingTiles(tarPos.Position)
				for _, tile := range tiles {
					d := w.Grid.GetDistance(tile, posComp.Position)
					if d < closestDistance {
						closestFreeTile = tile
						closestDistance = d
					}
				}
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

	return MovementAction{Target: closestFreeTile, Index: index}
}
