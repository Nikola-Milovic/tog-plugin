package systems

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/grid"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type TargetingSystem struct {
	World            *game.World
	WeightedEntities WeightedEntities // Who targets who
}

type weightedEntity struct {
	id     int
	weight int
}

// An WeightedEntities is a min-heap of entities with weight placed on how desirable of a target they are.
type WeightedEntities []weightedEntity

func (ts TargetingSystem) Update() {
	world := ts.World

	if world.Tick == 0 {
		ts.WeightedEntities = make([]weightedEntity, 0, 30)
		ts.startOfMatch()
		return
	}

	indexMap := world.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]
	attackComponents := world.ObjectPool.Components["AttackComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		movComp := movementComponents[index].(components.MovementComponent)
		posComp := positionComponents[index].(components.PositionComponent)
		atkComp := attackComponents[index].(components.AttackComponent)

		//chooseNewTarget := false
		//if !entities[indexMap[atkComp.Target]].Active {
		//	chooseNewTarget = true
		//}

		tag := entities[index].PlayerTag
		entId := entities[index].ID

		//	attackRange := atkComp.Range

		engagingDistance := movComp.MovementSpeed*20 + atkComp.Range

		tId := atkComp.Target
		tIndex := indexMap[tId]
		tPos := world.ObjectPool.Components["PositionComponent"][tIndex].(components.PositionComponent).Position
		//	distance := math.GetDistance(tPos, posComp.Position)

		//Move towards target goal
		switch entities[index].State {
		case constants.StateWalking: // keep moving towards target
			{
				movComp.Goal = tPos
			}
		case constants.StateEngaging: // choose an empty spot around our target
			{

				if world.Tick%5 != 0 {
					found, id := ts.selectTarget(entId, tag, math.Square(posComp.Position, posComp.Radius+atkComp.Range + 32))
					if found {
						if id != tId {
							blackboardIndex := findElement(world.Blackboard[tId], entId)                // find us in the previous targets blackboard
							world.Blackboard[tId] = removeIndex(world.Blackboard[tId], blackboardIndex) // remove us from the previous one
							world.Blackboard[id] = append(world.Blackboard[id], entId)                  // add us to the new target
							atkComp.Target = id
							movComp.Goal = positionComponents[indexMap[tId]].(components.PositionComponent).Position
						}
					} else {
						movComp.Goal = tPos
					}

				} else {
					found, id := ts.selectTarget(entId, tag, math.Square(posComp.Position, engagingDistance*0.3))

					if found {
						if id != tId {
							blackboardIndex := findElement(world.Blackboard[tId], entId)                // find us in the previous targets blackboard
							world.Blackboard[tId] = removeIndex(world.Blackboard[tId], blackboardIndex) // remove us from the previous one
							world.Blackboard[id] = append(world.Blackboard[id], entId)                  // add us to the new target
							atkComp.Target = id
							movComp.Goal = positionComponents[indexMap[tId]].(components.PositionComponent).Position
						}
					} else {
						movComp.Goal = tPos
					}
				}


			}
		case constants.StateThinking: // find enemy anywhere on map
			{
				found, id := ts.selectTarget(entId, tag, math.Square(posComp.Position, 200))

				if found {
					if id != tId {
						blackboardIndex := findElement(world.Blackboard[tId], entId)                // find us in the previous targets blackboard
						world.Blackboard[tId] = removeIndex(world.Blackboard[tId], blackboardIndex) // remove us from the previous one
						world.Blackboard[id] = append(world.Blackboard[id], entId)                  // add us to the new target
						atkComp.Target = id
						movComp.Goal = tPos
					} else {
						movComp.Goal = tPos
					}
				}
			}
		}

		world.ObjectPool.Components["AttackComponent"][index] = atkComp
		//	world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movComp
	}
}

func (ts TargetingSystem) startOfMatch() { // find a target for each unit
	p0Square := math.Square(math.Vector{X: float32(constants.MapWidth - constants.MapWidth/4), Y: float32(constants.MapHeight / 2)}, constants.MapWidth/4)
	p1Square := math.Square(math.Vector{X: float32(constants.MapWidth / 4), Y: float32(constants.MapHeight / 2)}, constants.MapWidth/4)

	world := ts.World
	indexMap := world.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]
	attackComponents := world.ObjectPool.Components["AttackComponent"]

	for index := range entities {
		movComp := movementComponents[index].(components.MovementComponent)
		posComp := positionComponents[index].(components.PositionComponent)
		atkComp := attackComponents[index].(components.AttackComponent)

		//	targetIndex := indexMap[movementComp.Goal]
		//targetPos := positionComponents[targetIndex].(components.PositionComponent)

		tag := entities[index].PlayerTag
		entID := entities[index].ID

		if tag == 0 {
			found, id := ts.selectTarget(entID, tag, p0Square)
			if found {
				world.Blackboard[id] = append(world.Blackboard[id], entID)
				tPos := positionComponents[indexMap[id]].(components.PositionComponent).Position
				movComp.Velocity = posComp.Position.To(tPos).Normalize()
				atkComp.Target = id
				movComp.Goal = tPos
				entities[index].State = constants.StateWalking
			}
		} else {
			found, id := ts.selectTarget(entID, tag, p1Square)
			if found {
				world.Blackboard[id] = append(world.Blackboard[id], entID)
				tPos := positionComponents[indexMap[id]].(components.PositionComponent).Position
				movComp.Velocity = posComp.Position.To(tPos).Normalize()
				atkComp.Target = id
				movComp.Goal = tPos
				entities[index].State = constants.StateWalking
			}
		}

		world.ObjectPool.Components["AttackComponent"][index] = atkComp
		world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movComp
	}

	//set velocities towards their target

}

//returns whether or not it found any viable target, if it did it returns the ID of the target
func (ts TargetingSystem) selectTarget(entityID, tag int, searchZone math.AABB) (bool, int) {

	ts.WeightedEntities = ts.WeightedEntities[:0]
	//Get all entities
	ts.World.Buff = ts.World.SpatialHash.Query(searchZone, ts.World.Buff[:0], helper.GetEnemyTag(tag))

	entities := ts.World.EntityManager.GetEntities()
	indexMap := ts.World.EntityManager.GetIndexMap()

	index := indexMap[entityID]

	movementComponents := ts.World.ObjectPool.Components["MovementComponent"]
	positionComponents := ts.World.ObjectPool.Components["PositionComponent"]
	attackComponents := ts.World.ObjectPool.Components["AttackComponent"]

	movComp := movementComponents[index].(components.MovementComponent)
	posComp := positionComponents[index].(components.PositionComponent)
	atkComp := attackComponents[index].(components.AttackComponent)

	for _, id := range ts.World.Buff {
		if id == entityID {
			continue
		}

		otherIndex := indexMap[id]
		otherEntity := entities[otherIndex]

		if otherEntity.PlayerTag == tag {
			fmt.Printf("We are the same tag")
			continue
		}

		otherPosComp := positionComponents[otherIndex].(components.PositionComponent)
		//otherMovComp := movementComponents[otherIndex].(components.MovementComponent)

		//	entityIndex := indexMap[id]
		//	entity := entities[entityIndex]
		dist := math.GetDistance(otherPosComp.Position, posComp.Position)
		w := int(dist)

		w *= int(ts.World.Grid.GetOccupationalMap().GetCell(grid.GlobalCordToTiled(otherPosComp.Position)))

		if dist <= atkComp.Range+otherPosComp.Radius+posComp.Radius+movComp.MovementSpeed*2 {
			w = w * 1 / 5
		}

		if dist <= atkComp.Range+otherPosComp.Radius+posComp.Radius {
			w = w* 1/10
		}
		//	x, y := grid.GlobalCordToTiled(otherPosComp.Position)

		//	w = w * 1/ int(ts.World.Grid.GetEnemyProximityImap(tag).GetCell(x, y))
		maxTars := ts.getMaxTargetsForEntity(id)
		tars := len(ts.World.Blackboard[id])
		if tars >= maxTars {
			w = w * 4
		} else {
			w = w * 2 / (5 + (maxTars - tars))
		}

		ts.WeightedEntities = append(ts.WeightedEntities, weightedEntity{weight: w, id: id})
	}

	if len(ts.WeightedEntities) > 0 {
		lowest := ts.WeightedEntities[0]
		for _, w := range ts.WeightedEntities {
			if w.weight < lowest.weight {
				lowest = w
			}
		}
		return true, lowest.id
	}

	return false, -1
}

func (ts TargetingSystem) getMaxTargetsForEntity(id int) int {
	index := ts.World.EntityManager.GetIndexMap()[id]
	//	movComp := ts.World.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)
	posComp := ts.World.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	//atkComp := ts.World.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)

	switch posComp.Radius {
	case 16:
		{
			return 3
		}
	case 20:
		{
			return 4
		}
	case 10:
		{
			return 2
		}
	}

	return 4
}

func findElement(s []int, el int) int {
	for i, e := range s {
		if e == el {
			return i
		}
	}
	return -1
}

func removeIndex(s []int, i int) []int {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}
