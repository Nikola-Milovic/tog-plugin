package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type SeparationSystem struct {
	World *game.World
	Buff  []int
}

func (ms SeparationSystem) Update() {
	world := ms.World
	entities := world.EntityManager.GetEntities()
//	movementComponents := world.ObjectPool.Components["MovementComponent"]
	attackComponents := world.ObjectPool.Components["AttackComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	if ms.Buff == nil {
		ms.Buff = make([]int, 100)
	}

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		posComp := positionComponents[index].(components.PositionComponent)
		atkComp := attackComponents[index].(components.AttackComponent)

		ms.Buff = helper.GetNearbyEntities(150, world, index, ms.Buff[:0])

		dX := float32(0)
		dY := float32(0)

		edX := float32(0)
		edY := float32(0)

		for _, ent := range ms.Buff {
			if ent == index {
				continue
			}
			entity := entities[ent]
			me := entities[index]

			myPos := positionComponents[index].(components.PositionComponent)
			otherPos := positionComponents[ent].(components.PositionComponent)

			if entity.PlayerTag == me.PlayerTag {//ally
				detectionLimit := myPos.BoundingBox.X/2 + otherPos.BoundingBox.X/2 + 4
				distance := math.GetDistance(myPos.Position, otherPos.Position)

				if distance < detectionLimit {
					dX -= otherPos.Position.X - myPos.Position.X
					dY -= otherPos.Position.Y - myPos.Position.Y
				}
			} else {
				detectionLimit := myPos.BoundingBox.X/2 + atkComp.Range
				distance := math.GetDistance(myPos.Position, otherPos.Position)

				if distance < detectionLimit {
					dX -= otherPos.Position.X - myPos.Position.X
					dY -= otherPos.Position.Y - myPos.Position.Y
				}
			}

		}

		posComp.Position.X += dX/5
		posComp.Position.Y += dY/5

		posComp.Position.X += edX
		posComp.Position.Y += edY

		world.ObjectPool.Components["PositionComponent"][index] = posComp

	}
}
