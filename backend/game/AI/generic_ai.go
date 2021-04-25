package ai

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type GenericAI struct {
	World *game.World
	Buff  []int
}

func (ai GenericAI) PerformAI(index int) {
	world := ai.World
	//g := ai.World.Grid
	indexMap := world.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
	atkComp := world.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	posComp := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	movComp := world.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)

	playerTag := entities[index].PlayerTag
	attackRange := atkComp.Range
	engagingDistance := movComp.MovementSpeed*20 + attackRange

	fmt.Println(entities[index].State)

	//	fmt.Printf("I %d am at state %s\n", index, entities[index].State)

	switch entities[index].State {
	case constants.StateWalking:
		{
			targetIndex := world.EntityManager.GetIndexMap()[movComp.TargetID]
			targetPos := world.ObjectPool.Components["PositionComponent"][targetIndex].(components.PositionComponent)

			dist := math.GetDistanceIncludingDiagonalVectors(targetPos.Position, posComp.Position) - posComp.Radius - targetPos.Radius
			if dist <= engagingDistance {
				helper.SwitchState(entities, index, constants.StateEngaging, world)
			}
		}
	case constants.StateEngaging:
		{ // check if there is a closer unit
			ai.Buff = world.SpatialHash.Query(math.Square(posComp.Position, engagingDistance), ai.Buff[:0], getEnemyTag(playerTag), true)

			closestEnemyIndex := -1
			closestEnemyDistance := float32(10000)

			for _, id := range ai.Buff {
				eIndex := indexMap[id]

				ePos := world.ObjectPool.Components["PositionComponent"][eIndex].(components.PositionComponent)
				dist := math.GetDistanceIncludingDiagonalVectors(ePos.Position, posComp.Position) - posComp.Radius - ePos.Radius

				if dist < closestEnemyDistance {
					closestEnemyDistance = dist
					closestEnemyIndex = eIndex
				}

				if dist <= attackRange {
					attackTarget(index, entities[eIndex].ID, entities[index].ID, world)
					return
				}

			}
			if closestEnemyIndex == -1 {
				return
			}

			if closestEnemyDistance <= engagingDistance {
				moveTowardsTarget(index, entities[closestEnemyIndex].ID, world, true)
			} else {
				moveTowardsTarget(index, entities[closestEnemyIndex].ID, world, false)
			}

		}

	case constants.StateAttacking:
		{
			return
		}

	case constants.StateThinking: // find closest enemy
		{
			closestEnemyIndex := -1
			closestEnemyDistance := float32(10000)
			for eIndx, p := range world.ObjectPool.Components["PositionComponent"] {
				if entities[eIndx].PlayerTag == playerTag {
					continue
				}
				pos := p.(components.PositionComponent)
				dist := math.GetDistanceIncludingDiagonalVectors(pos.Position, posComp.Position) - posComp.Radius - pos.Radius
				if dist <= attackRange {
					attackTarget(index, entities[eIndx].ID, entities[index].ID, world)
					return
				}

				if dist < closestEnemyDistance {
					closestEnemyDistance = dist
					closestEnemyIndex = eIndx
				}
			}
			if closestEnemyIndex == -1 {
				return
			}

			if closestEnemyDistance <= engagingDistance {
				moveTowardsTarget(index, entities[closestEnemyIndex].ID, world, true)
			} else {
				moveTowardsTarget(index, entities[closestEnemyIndex].ID, world, false)
			}
		}
	}
}
