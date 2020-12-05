package game

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type KnightAI struct {
	world *World
}

func (ai KnightAI) CalculateAction(index int) engine.Action {
	w := ai.world

	nearbyEntities := GetNearbyEntities(30, w, index)
	atkComp := w.ObjectPool.Components["AttackComponent"][index].(AttackComponent)
	posComp := w.ObjectPool.Components["PositionComponent"][index].(PositionComponent)

	//If we're already attacking, keep attacking
	if atkComp.Target != -1 {
		tarPos := w.ObjectPool.Components["PositionComponent"][atkComp.Target].(PositionComponent)
		if w.Grid.GetDistanceIncludingDiagonal(posComp.Position, tarPos.Position) < 2 {
			return AttackAction{Target: atkComp.Target}
		}
	}

	//Check if an enemy is in range or move to somewhere
	closestFreeTile := engine.Vector{}
	closestDistance := 100000
	for _, indx := range nearbyEntities {
		if w.EntityManager.Entities[index].PlayerTag != w.EntityManager.Entities[indx].PlayerTag {
			tarPos := w.ObjectPool.Components["PositionComponent"][indx].(PositionComponent)
			if w.Grid.GetDistanceIncludingDiagonal(tarPos.Position, posComp.Position) < 2 {
				return AttackAction{Target: atkComp.Target}
			}

			tiles := w.Grid.GetSurroundingTiles(posComp.Position.X, posComp.Position.Y)
			for _, tile := range tiles {
				if w.Grid.GetDistance(tile, posComp.Position) < closestDistance {
					closestFreeTile = tile
				}
			}

		}
	}

	//Reset target to noone
	atkComp.Target = -1
	w.ObjectPool.Components["AttackComponent"][index] = atkComp

	return MovementAction{Target: closestFreeTile}
}
