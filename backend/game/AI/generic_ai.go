package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type GenericAI struct {
	World *game.World
	Buff  []int
}

func (ai GenericAI) PerformAI(index int) {
	world := ai.World
	//g := ai.World.Grid
	//	indexMap := world.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
	atkComp := world.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	//posComp := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	movComp := world.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)

	unitID := entities[index].ID

	playerTag := entities[index].PlayerTag
	attackRange := atkComp.Range
	engagingDistance := movComp.MovementSpeed*20 + attackRange

	switch entities[index].State {
	case constants.StateWalking:
		{
			tId, dist := findClosestEnemy(ai.Buff, playerTag, index, 300, world)
			if tId == -1 {

			} else {
				if dist <= engagingDistance {
					moveTowardsTarget(index, tId, world, true)
				}
			}
		}
	case constants.StateEngaging:
		{
			tId, dist := findClosestEnemy(ai.Buff, playerTag, index, engagingDistance, world)
			if tId == -1 {

			} else {
				if dist <= attackRange {
					attackTarget(index, tId, unitID, world)
				} else if dist > engagingDistance {
					moveTowardsTarget(index, tId, world, false)
				}
			}
		}
	case constants.StateAttacking:
		{

		}
	case constants.StateThinking:
		{
			tId, dist := findClosestEnemy(ai.Buff, playerTag, index, 800, world)
			if tId == -1 {
			} else {
				if dist <= attackRange {
					attackTarget(index, tId, unitID, world)
				} else if dist <= engagingDistance {
					moveTowardsTarget(index, tId, world, true)
				} else {
					moveTowardsTarget(index, tId, world, false)
				}
			}
		}
	}
}

func findClosestEnemy(buff []int, unitTag int, index int, searchRadius float32, world *game.World) (id int, dist float32) {
	indexMap := world.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
	posComp := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)

	buff = world.SpatialHash.Query(math.Square(posComp.Position, searchRadius), buff[:0], getEnemyTag(unitTag), true)

	closestEnemyIndex := -1
	closestEnemyDistance := float32(10000)

	for _, i := range buff {
		eIndex := indexMap[i]

		ePos := world.ObjectPool.Components["PositionComponent"][eIndex].(components.PositionComponent)
		d := math.GetDistanceIncludingDiagonalVectors(ePos.Position, posComp.Position) - posComp.Radius - ePos.Radius

		if d < closestEnemyDistance {
			closestEnemyDistance = d
			closestEnemyIndex = eIndex
		}
	}
	if closestEnemyIndex == -1 {
		return -1, -1
	}

	return entities[closestEnemyIndex].ID, closestEnemyDistance
}
