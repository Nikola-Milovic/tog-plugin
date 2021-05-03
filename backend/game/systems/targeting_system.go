package systems

import (
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

	//g := world.Grid

	indexMap := world.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	//positionComponents := world.ObjectPool.Components["PositionComponent"]
	attackComponents := world.ObjectPool.Components["AttackComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		movementComp := movementComponents[index].(components.MovementComponent)
		//posComp := positionComponents[index].(components.PositionComponent)
		atkComp := attackComponents[index].(components.AttackComponent)

		//	targetIndex := indexMap[movementComp.Goal]
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
					tIndex := indexMap[tId]
					tPos := world.ObjectPool.Components["PositionComponent"][tIndex].(components.PositionComponent).Position
					if dist <= engagingDistance {
						helper.MoveTowardsTarget(index, tPos, world, true)
					}
				}
			}
		case constants.StateEngaging:
			{
				tId, dist := helper.FindClosestEnemy(world.Buff, tag, index, engagingDistance, world)
				if tId == -1 {
				} else {
					tIndex := indexMap[tId]
					tPos := world.ObjectPool.Components["PositionComponent"][tIndex].(components.PositionComponent).Position
					if dist < attackRange {
						//	fmt.Printf("Distance to target is %.2f, will attack range is %.2f\n", dist, attackRange)
						helper.AttackTarget(index, tId, unitID, world)
					} else if dist > engagingDistance {
						helper.MoveTowardsTarget(index, tPos, world, false)
					} else {
						//size := int((64+posComp.Radius+attackRange)/constants.TileSize + 1)
						//
						//ex, ey := grid.GlobalCordToTiled(tPos)
						//mx, my := grid.GlobalCordToTiled(posComp.Position)
						//wm := g.GetWorkingMap(size, size)
						//
						//engine.AddIntoSmallerMap(g.GetEnemyProximityImap(tag), wm, ex, ey, 1)
						//engine.AddIntoSmallerMap(g.GetProximityImaps()[tag], wm, ex, ey, -0.6)
						//engine.AddIntoBiggerMap(grid.GetProximityTemplate(movementComp.MovementSpeed).Imap, wm, math.Clamp(ex+(ex-mx), 0, wm.Width), math.Clamp(ey+(ey-my), 0, wm.Height), 0.6)
						////	engine.AddIntoSmallerMap(g.GetOccupationalMap(), wm, ex, ey, -1)
						//engine.AddIntoSmallerMap(g.GetGoalMap(), wm, ex, ey, -1)
						//engine.AddIntoBiggerMap(g.GetInterestTemplate(10), wm, size/2, size/2, 2)
						//
						//x, y, _ := wm.GetHighestCell()
						//adjX, adjY := grid.GetBaseMapCoordsFromSectionImapCoords(ex, ey, x, y)
						//
						//engine.AddIntoBiggerMap(grid.GetSizeTemplate(posComp.Radius).Imap, g.GetGoalMap(), adjX, adjY, 1)
						//
						//if g.GetOccupationalMap().GetCell(adjX, adjY) > 0 {
						//	fmt.Printf("Goal cell occupied\n")
						//	helper.SwitchState(entities, index, constants.StateThinking, world)
						//}
						//
						movementComp.Goal = tPos //math.Vector{X: float32(adjX * constants.TileSize), Y: float32(adjY * constants.TileSize)}
						movementComponents[index] = movementComp
					}
				}
			}
		case constants.StateAttacking:
			{

			}
		}
	}
}
