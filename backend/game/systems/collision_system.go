package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type CollisionSystem struct {
	World *game.World
	Buff  []int
}

func (ms CollisionSystem) Update() {
	ms.overlapping()
	//	ms.collisionAvoidance()
}

func (ms CollisionSystem) overlapping() {
	world := ms.World
	indexMap := ms.World.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	//attackComponents := world.ObjectPool.Components["AttackComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		if entities[index].State == constants.StateAttacking {
			return
		}

		posComp := positionComponents[index].(components.PositionComponent)
		//atkComp := attackComponents[index].(components.AttackComponent)
		movComp := movementComponents[index].(components.MovementComponent)

		ms.Buff = world.SpatialHash.Query(math.Square(posComp.Position, posComp.Radius+80), ms.Buff[:0], -1, true)

		dX := float32(0)
		dY := float32(0)

		edX := float32(0)
		edY := float32(0)

		movComp.DestinationMultiplier = 1

		for _, otherID := range ms.Buff {
			otherIndex := indexMap[otherID]
			entity := entities[otherIndex]
			me := entities[index]

			if me.ID == otherID {
				continue
			}

			otherPos := positionComponents[otherIndex].(components.PositionComponent)
			//	otherMovComp := movementComponents[otherIndex].(components.MovementComponent)

			//Prevent overlapping
			if entity.PlayerTag == me.PlayerTag {
				detectionLimit := posComp.Radius/2 + otherPos.Radius
				distance := math.GetDistance(posComp.Position, otherPos.Position)

				//Overlap
				if distance < detectionLimit && distance > 5 { // already colliding
					dX -= otherPos.Position.X - posComp.Position.X
					dY -= otherPos.Position.Y - posComp.Position.Y
				} else if distance < 5 {
					dX -= 5 * (otherPos.Position.X - posComp.Position.X)
					dY -= 5 * (otherPos.Position.Y - posComp.Position.Y)
					if dX == 0 {
						dX = 10
					}
					if dY == 0 {
						dY = 10
					}
				}

				if distance < detectionLimit {
					diff := detectionLimit - distance
					adjustment := otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(diff)
					movComp.Velocity = movComp.Velocity.Subtract(adjustment).Normalize().MultiplyScalar(movComp.MovementSpeed)
					movComp.DestinationMultiplier = 0
				}

				// Seperation force
				//if entity.State == constants.StateAttacking {
				//	forceFieldRadius := 2*detectionLimit + 12
				//	if distance > detectionLimit && distance < forceFieldRadius {
				//		diff := forceFieldRadius - distance
				//		adjustment := otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(diff)
				//		movComp.Velocity = movComp.Velocity.Subtract(adjustment).Normalize().MultiplyScalar(movComp.MovementSpeed)
				//		movComp.DestinationMultiplier = 0
				//	}
				//}
			}
		}

		posComp.Position.X += edX
		posComp.Position.Y += edY

		posComp.Position.X += dX / 4
		posComp.Position.Y += dY / 4

		world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movComp
	}
}

func (ms CollisionSystem) collisionAvoidance() {
	world := ms.World
	indexMap := ms.World.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	//	attackComponents := world.ObjectPool.Components["AttackComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		if entities[index].State == constants.StateAttacking {
			continue
		}

		unitTag := entities[index].PlayerTag

		posComp := positionComponents[index].(components.PositionComponent)
		//	atkComp := attackComponents[index].(components.AttackComponent)
		movComp := movementComponents[index].(components.MovementComponent)

		velocity := movComp.Velocity

		futurePosition := posComp.Position.Multiply(movComp.Velocity.MultiplyScalar(3))

		ms.Buff = world.SpatialHash.Query(math.Square(futurePosition, posComp.Radius+120), ms.Buff[:0], unitTag, true)

		for _, ent := range ms.Buff {
			eIndex := indexMap[ent]
			if eIndex == index {
				continue
			}
			other := entities[eIndex]
			//me := entities[index]

			otherPos := positionComponents[eIndex].(components.PositionComponent)
			//otherMovComp := movementComponents[ent].(components.MovementComponent)

			detectionLimit := posComp.Radius + otherPos.Radius + 32
			distance := math.GetDistance(futurePosition, otherPos.Position)

			if distance < detectionLimit { // will collide, avoid
				if other.State == constants.StateAttacking {
					diff := detectionLimit - distance
					adj := otherPos.Position.To(futurePosition).Normalize().MultiplyScalar(diff * 2)

					velocity = velocity.Subtract(adj)
				} else {
					diff := detectionLimit - distance
					adj := otherPos.Position.To(futurePosition).Normalize().MultiplyScalar(diff)

					velocity = velocity.Subtract(adj)
				}
			}
		}

		movComp.Velocity = velocity

		world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movComp
	}
}
