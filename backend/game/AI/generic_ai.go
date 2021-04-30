package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
)

type GenericAI struct {
	World *game.World
}

func (ai GenericAI) PerformAI(index int) {
	world := ai.World
	//g := ai.World.Grid
	//	indexMap := world.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
	atkComp := world.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	//posComp := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	movComp := world.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)

	if entities[index].State != constants.StateThinking {
		return
	}

	unitID := entities[index].ID

	playerTag := entities[index].PlayerTag
	attackRange := atkComp.Range
	engagingDistance := movComp.MovementSpeed*20 + attackRange

	switch entities[index].State {
	case constants.StateWalking:
		{

		}
	case constants.StateEngaging:
		{

		}
	case constants.StateAttacking:
		{

		}
	case constants.StateThinking:
		{
			tId, dist := helper.FindClosestEnemy(world.Buff, playerTag, index, 800, world)
			if tId == -1 {
			} else {
				if dist <= attackRange {
					helper.AttackTarget(index, tId, unitID, world)
				} else if dist <= engagingDistance {
					helper.MoveTowardsTarget(index, tId, world, true)
				} else {
					helper.MoveTowardsTarget(index, tId, world, false)
				}
			}
		}
	}
}