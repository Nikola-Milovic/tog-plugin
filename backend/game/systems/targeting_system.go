package systems

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
)

//TargetingSystem should
//1) Choose the closest enemy
//2) Check whether or not our current target is the best one
//3) If engaging an enemy, find a free spot for us to navigate to
type TargetingSystem struct {
	World *game.World
}

func (ts TargetingSystem) Update() {
	world := ts.World

	if world.Tick%5 != 0 {
		return
	}

	//	indexMap := World.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	//	positionComponents := World.ObjectPool.Components["PositionComponent"]
	attackComponents := world.ObjectPool.Components["AttackComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		movementComp := movementComponents[index].(components.MovementComponent)
		//	posComp := positionComponents[index].(components.PositionComponent)
		atkComp := attackComponents[index].(components.AttackComponent)

		//	targetIndex := indexMap[movementComp.TargetID]
		//targetPos := positionComponents[targetIndex].(components.PositionComponent)

		tag := entities[index].PlayerTag
		unitID := entities[index].ID

		attackRange := atkComp.Range

		engagingDistance := movementComp.MovementSpeed*20 + attackRange

		switch entities[index].State {
		case constants.StateWalking:
			{
				tId, dist := helper.FindClosestEnemy(world.Buff, tag, index, 300, world)
				if tId == -1 {

				} else {
					if dist <= engagingDistance {
						helper.MoveTowardsTarget(index, tId, world, true)
					}
				}
			}
		case constants.StateEngaging:
			{
				tId, dist := helper.FindClosestEnemy(world.Buff, tag, index, engagingDistance, world)
				if tId == -1 {
				} else {
					if dist < attackRange {
						fmt.Printf("Distance to target is %.2f, will attack range is %.2f\n", dist, attackRange)
						helper.AttackTarget(index, tId, unitID, world)
					} else if dist > engagingDistance {
						helper.MoveTowardsTarget(index, tId, world, false)
					}
				}
			}
		case constants.StateAttacking:
			{

			}
		}
	}
}
