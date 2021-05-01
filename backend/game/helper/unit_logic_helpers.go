package helper

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

func AttackTarget(index int, target int, id int, w *game.World) {
	attackComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	attackComp.Target = target
	attackComp.TimeSinceLastAttack = w.Tick
	w.ObjectPool.Components["AttackComponent"][index] = attackComp

	SwitchState(w.EntityManager.GetEntities(), index, constants.StateAttacking, w)

	clientEvent := make(map[string]interface{}, 3)
	clientEvent["event"] = "attack"
	clientEvent["me"] = id
	clientEvent["who"] = target
	w.ClientEventManager.AddEvent(clientEvent)
}

func MoveTowardsTarget(index int, tarPos math.Vector, w *game.World, isEngaging bool) {
	movementComp := w.GetObjectPool().Components["MovementComponent"][index].(components.MovementComponent)

	if isEngaging {
		SwitchState(w.EntityManager.GetEntities(), index, constants.StateEngaging, w)
	} else {
		SwitchState(w.EntityManager.GetEntities(), index, constants.StateWalking, w)
	}
	movementComp.Goal = tarPos

	w.GetObjectPool().Components["MovementComponent"][index] = movementComp
}

func GetEnemyTag(tag int) int {
	if tag == 0 {
		return 1
	} else {
		return 0
	}
}

func FindClosestEnemy(buff []int, unitTag int, index int, searchRadius float32, world *game.World) (id int, dist float32) {
	indexMap := world.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
	posComp := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)

	buff = world.SpatialHash.Query(math.Square(posComp.Position, searchRadius), buff[:0], GetEnemyTag(unitTag), true)

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

