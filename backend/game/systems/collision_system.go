package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type SeparationSystem struct {
	World *game.World
	Buff  []int
}

func (ms SeparationSystem) Update() {
	world := ms.World
	indexMap := ms.World.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
//	movementComponents := world.ObjectPool.Components["MovementComponent"]
	attackComponents := world.ObjectPool.Components["AttackComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		if entities[index].State == constants.StateAttacking {
			continue
		}

		posComp := positionComponents[index].(components.PositionComponent)
		atkComp := attackComponents[index].(components.AttackComponent)
		//	myMovComp := movementComponents[index].(components.MovementComponent)

		ms.Buff = world.SpatialHash.Query(math.Square(posComp.Position, 120), ms.Buff[:0], -1, true)

		dX := float32(0)
		dY := float32(0)

		edX := float32(0)
		edY := float32(0)

		for _, ent := range ms.Buff {
			eIndex := indexMap[ent]
			if ent == index {
				continue
			}
			entity := entities[eIndex]
			me := entities[index]

			otherPos := positionComponents[eIndex].(components.PositionComponent)

			//otherMovComp := movementComponents[ent].(components.MovementComponent)

			if entity.PlayerTag == me.PlayerTag { //ally
				detectionLimit := posComp.Radius + otherPos.Radius - 4
				distance := math.GetDistance(posComp.Position, otherPos.Position)

				if distance < detectionLimit { // already colliding
					dX -= otherPos.Position.X - posComp.Position.X
					dY -= otherPos.Position.Y - posComp.Position.Y
				}

			} else {
				detectionLimit := posComp.Radius + atkComp.Range - 4
				distance := math.GetDistance(posComp.Position, otherPos.Position)

				if distance < detectionLimit {
					dX -= otherPos.Position.X - posComp.Position.X
					dY -= otherPos.Position.Y - posComp.Position.Y
				}
			}

		}

		posComp.Position.X += dX / 5
		posComp.Position.Y += dY / 5

		posComp.Position.X += edX
		posComp.Position.Y += edY

		world.ObjectPool.Components["PositionComponent"][index] = posComp
		//world.ObjectPool.Components["MovementComponent"][index] = myMovComp

	}
}
