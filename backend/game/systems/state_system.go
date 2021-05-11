package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

//StateSystem should
//1) Choose the closest enemy
//2) Check whether or not our current target is the best one
//3) If engaging an enemy, find a free spot for us to navigate to
type StateSystem struct {
	World *game.World
}

func (ss StateSystem) Update() {
	world := ss.World

	indexMap := world.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]
	attackComponents := world.ObjectPool.Components["AttackComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		movementComp := movementComponents[index].(components.MovementComponent)
		posComp := positionComponents[index].(components.PositionComponent)
		atkComp := attackComponents[index].(components.AttackComponent)

		//	targetIndex := indexMap[movementComp.Goal]
		//targetPos := positionComponents[targetIndex].(components.PositionComponent)

		//	tag := entities[index].PlayerTag
		unitID := entities[index].ID

		attackRange := atkComp.Range

		engagingDistance := movementComp.MovementSpeed*20 + attackRange

		tId := atkComp.Target
		tIndex := indexMap[tId]
		tarEntity := entities[tIndex]

		if !tarEntity.Active {
			entities[index].State = constants.StateThinking
			continue
		}

		tPos := world.ObjectPool.Components["PositionComponent"][tIndex].(components.PositionComponent).Position
		tRadius := world.ObjectPool.Components["PositionComponent"][tIndex].(components.PositionComponent).Radius
		distance := math.GetDistance(tPos, posComp.Position) - posComp.Radius -  tRadius

		switch entities[index].State {
		case constants.StateWalking:
			{
				if distance <= engagingDistance {
					helper.MoveTowardsTarget(index, tPos, world, true)
				}
			}
		case constants.StateEngaging:
			{
				if distance-1 <= attackRange {
					//	fmt.Printf("Distance to target is %.2f, will attack range is %.2f\n", dist, attackRange)
					helper.AttackTarget(index, tId, unitID, world)
				} else if distance > engagingDistance {
					helper.MoveTowardsTarget(index, tPos, world, false)
				}
			}
		case constants.StateAttacking:
			{
			}
		case constants.StateThinking: {
			if distance <= attackRange {
				//	fmt.Printf("Distance to target is %.2f, will attack range is %.2f\n", dist, attackRange)
				helper.AttackTarget(index, tId, unitID, world)
			} else if distance <= engagingDistance {
				helper.MoveTowardsTarget(index, tPos, world, true)
			} else {
				helper.MoveTowardsTarget(index, tPos, world, false)
			}
		}
		}
	}
}
