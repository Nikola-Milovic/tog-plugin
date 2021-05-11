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
		movComp := movementComponents[index].(components.MovementComponent)

		world.Buff = world.SpatialHash.Query(math.Square(posComp.Position, posComp.Radius+80), world.Buff[:0], -1)

		goal := movComp.Goal

	//	distToGoal := math.GetDistance(goal, posComp.Position)

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

				adjustment := math.Zero()

				if distance < detectionLimit-4 { // overlapping
					adjustment = otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(2)
					movComp.Acceleration = movComp.Acceleration.Add(adjustment)
					continue
				}

				if otherMovComp.Velocity != math.Zero() { // seperation while walking
					if distance < detectionLimit*1.5{
						adjustment = otherPos.Position.To(posComp.Position).Normalize()
						if otherPos.Radius == posComp.Radius { // we are equal, check speed
							if otherMovComp.MovementSpeed > movComp.MovementSpeed { // they are faster, we should just slow down a bit
								adjustment = adjustment.MultiplyScalar(0.4)
							} else if otherMovComp.MovementSpeed < movComp.MovementSpeed { // we are faster they should move
								adjustment =  adjustment.MultiplyScalar(0.3)
							} else { // we are equal, just seperate a bit
								adjustment = adjustment.MultiplyScalar(0.3)
							}
						} else if otherPos.Radius < posComp.Radius { // we are bigger, they move

						} else { // they are bigger we go around

						}
					}

					a := math.Abs(otherMovComp.Velocity.AngleTo(movComp.Velocity))
					if a < 10 {
						adjustment = adjustment.MultiplyScalar(0.7)
						if a == 0 {
							adjustment = adjustment.Add(movComp.Velocity.PerpendicularClockwise().Normalize().MultiplyScalar(0.2))
						}
					}

				} else {
					if distance < detectionLimit { // other is standing
						adjustment = otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(2)
						adjustment = adjustment.Add(posComp.Position.To(goal).Normalize().MultiplyScalar(2))
					}
				}

				movComp.Acceleration = movComp.Acceleration.Add(adjustment)
			} else {
				detectionLimit := posComp.Radius + otherPos.Radius + atkComp.Range + 2
				distance := math.GetDistance(posComp.Position, otherPos.Position)

				//Overlap
				if distance < detectionLimit*1.2 {
					adjustment := otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(0.7)
					movComp.Acceleration = movComp.Acceleration.Add(adjustment)
				}

			}
		}

		//if distToGoal < atkComp.Range + posComp.Radius + movComp.MovementSpeed*3 + 16 {
		//	movComp.Acceleration = movComp.Acceleration.MultiplyScalar(0.3)
		//}

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

		world.Buff = world.SpatialHash.Query(math.Square(futurePosition, posComp.Radius+60), world.Buff[:0], unitTag)

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

			distance := math.GetDistance(futurePosition, otherPos.Position.Add(otherMovComp.Velocity.MultiplyScalar(1.5)))

			adjustment := math.Zero()

			if distance < detectionLimit-4 { // overlapping
				adjustment = otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(1.2)
				movComp.Acceleration = movComp.Acceleration.Add(adjustment)
				continue
			}

			if otherMovComp.Velocity != math.Zero() { // seperation while walking
				if distance < detectionLimit*1.5 {
					adjustment = otherPos.Position.To(futurePosition).Normalize()
					if otherPos.Radius == posComp.Radius { // we are equal, check speed
						if otherMovComp.MovementSpeed > movComp.MovementSpeed { // they are faster, we should just slow down a bit
							adjustment = adjustment.MultiplyScalar(0.2)
						} else if otherMovComp.MovementSpeed < movComp.MovementSpeed { // we are faster they should move
							adjustment =  adjustment.MultiplyScalar(0.1)
						} else { // we are equal, just seperate a bit
							adjustment = adjustment.MultiplyScalar(0.1)
						}
					} else if otherPos.Radius < posComp.Radius { // we are bigger, they move

					} else { // they are bigger we go around

					}
				}

				a := math.Abs(otherMovComp.Velocity.AngleTo(movComp.Velocity))
				if a < 10 {
					adjustment = adjustment.MultiplyScalar(0.5)
					if a == 0 {
						adjustment = adjustment.Add(movComp.Velocity.PerpendicularClockwise().Normalize().MultiplyScalar(0.1))
					}
				}

			} else {
				if distance < detectionLimit {
					adjustment = otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(0.8)
				}
			}

			movComp.Acceleration = movComp.Acceleration.Add(adjustment)

			world.ObjectPool.Components["MovementComponent"][eIndex] = otherMovComp
		}

		movComp.Velocity = velocity

		world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movComp
	}
}
