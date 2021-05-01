package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

//CollisionSystem is supposed to do
//1) Not allow overlapping of units
//2) Keep units from bumping into each other
//3) Let some units pass before others/ go around someone
type CollisionSystem struct {
	World *game.World
}

func (ms CollisionSystem) Update() {
	ms.collision()
	//ms.collisionPrevention()
}

func (ms CollisionSystem) collision() {
	world := ms.World

	if world.Tick == 0 {
		return
	}

	indexMap := ms.World.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	//attackComponents := World.ObjectPool.Components["AttackComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		if entities[index].State == constants.StateAttacking {
			continue
		}

		posComp := positionComponents[index].(components.PositionComponent)
		//atkComp := attackComponents[index].(components.AttackComponent)
		movComp := movementComponents[index].(components.MovementComponent)

		world.Buff = world.SpatialHash.Query(math.Square(posComp.Position, posComp.Radius+80), world.Buff[:0], -1, true)

	//	targetPos := movComp.Goal

		dX := float32(0)
		dY := float32(0)

		edX := float32(0)
		edY := float32(0)

		for _, otherID := range world.Buff {
			otherIndex := indexMap[otherID]
			entity := entities[otherIndex]
			me := entities[index]

			if me.ID == otherID {
				continue
			}

			otherPos := positionComponents[otherIndex].(components.PositionComponent)
			otherMovComp := movementComponents[otherIndex].(components.MovementComponent)

			if entity.PlayerTag == me.PlayerTag {
				detectionLimit := posComp.Radius + otherPos.Radius
				distance := math.GetDistance(posComp.Position, otherPos.Position)

				//Overlap
				if distance < detectionLimit/2 {
					dX -= otherPos.Position.X - posComp.Position.X
					dY -= otherPos.Position.Y - posComp.Position.Y
				}

				if otherMovComp.Velocity != math.Zero() { // seperation while walking
					if distance < detectionLimit*1.5 {
						if otherPos.Radius == posComp.Radius { // we are equal, check speed
							if otherMovComp.MovementSpeed > movComp.MovementSpeed { // they are faster, we should just slow down a bit
								movComp.DestinationMultiplier = 0.6
								adjustment := otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(1.5)
								movComp.Velocity = movComp.Velocity.Add(adjustment)
								continue
							} else if otherMovComp.MovementSpeed < movComp.MovementSpeed { // we are faster they should move
								adjustment := otherPos.Position.To(posComp.Position).Normalize()
								movComp.Velocity = movComp.Velocity.Add(adjustment.MultiplyScalar(0.3))
								continue
							} else { // we are equal, just seperate a bit
								movComp.DestinationMultiplier = 0.8
								adjustment := otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(0.5)
								movComp.Velocity = movComp.Velocity.Add(adjustment)
								continue
							}
						} else if otherPos.Radius < posComp.Radius { // we are bigger, they move

						} else { // they are bigger we go around

						}
					}
				}
			} else {
				detectionLimit := posComp.Radius + otherPos.Radius - 4
				distance := math.GetDistance(posComp.Position, otherPos.Position)

				//Overlap
				if distance < detectionLimit {
					edX -= otherPos.Position.X - posComp.Position.X
					edY -= otherPos.Position.Y - posComp.Position.Y
				}
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

func (ms CollisionSystem) collisionPrevention() {
	world := ms.World
	indexMap := ms.World.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	//	attackComponents := World.ObjectPool.Components["AttackComponent"]
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

		futurePosition := posComp.Position.Add(movComp.Velocity.MultiplyScalar(3))

		world.Buff = world.SpatialHash.Query(math.Square(futurePosition, posComp.Radius+80), world.Buff[:0], unitTag, true)

		for _, ent := range world.Buff {
			eIndex := indexMap[ent]
			if eIndex == index {
				continue
			}
			//other := entities[eIndex]
			//me := entities[index]

			otherPos := positionComponents[eIndex].(components.PositionComponent)
			otherMovComp := movementComponents[eIndex].(components.MovementComponent)

			detectionLimit := posComp.Radius + otherPos.Radius + 16

			distance := math.GetDistance(futurePosition, otherPos.Position.Add(otherMovComp.Velocity.MultiplyScalar(3)))

			if distance < detectionLimit*1.5 && distance > detectionLimit-8 { // will collide, avoid
				if otherMovComp.Velocity != math.Zero() { // seperation while walking
					if otherPos.Radius == posComp.Radius { // we are equal, check speed
						if otherMovComp.MovementSpeed > movComp.MovementSpeed { // they are faster
							movComp.DestinationMultiplier = 0.6
							adjustment := otherPos.Position.To(futurePosition).Normalize().MultiplyScalar(0.2)
							movComp.Velocity = movComp.Velocity.Add(adjustment)
							continue
						} else if otherMovComp.MovementSpeed < movComp.MovementSpeed { // we are faster they should move
							adjustment := otherPos.Position.To(futurePosition).Normalize()
							movComp.Velocity = movComp.Velocity.Add(adjustment.MultiplyScalar(0.1))

							otheradj := futurePosition.To(otherPos.Position).Normalize().MultiplyScalar(2)
							otheradj = otheradj.Add(otherMovComp.Velocity.Cross().Normalize().MultiplyScalar(2))
							otherMovComp.Velocity = otherMovComp.Velocity.Add(otheradj)
							otherMovComp.DestinationMultiplier = 0.4

							continue
						} else { // we are equal, just seperate a bit
							//movComp.DestinationMultiplier = 0.8
							//adjustment := otherPos.Position.To(futurePosition).Normalize().MultiplyScalar(0.1)
							//movComp.Velocity = movComp.Velocity.Add(adjustment)
							continue
						}
					} else if otherPos.Radius < posComp.Radius { // we are bigger, they move

					} else { // they are bigger we go around

					}

				}
			}
			world.ObjectPool.Components["MovementComponent"][eIndex] = otherMovComp
		}

		movComp.Velocity = velocity

		world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movComp
	}
}
