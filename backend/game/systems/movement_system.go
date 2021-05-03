package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type MovementSystem struct {
	World *game.World
}

func (ms MovementSystem) Update() {
	world := ms.World
	//useless := 0
	//g := World.Grid
	indexMap := world.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		if entities[index].State != constants.StateWalking && entities[index].State != constants.StateEngaging {
			continue
		}

		//Setup
		playerTag := entities[index].PlayerTag

		movementComp := movementComponents[index].(components.MovementComponent)
		posComp := positionComponents[index].(components.PositionComponent)
		velocity := movementComp.Velocity

		targetPos := movementComp.Goal
		targetDir := posComp.Position.To(targetPos).Normalize()

		if entities[index].State == constants.StateEngaging {
			forwardRay := targetDir.Normalize()
			//rightRay := forwardRay.PerpendicularClockwise()
			//leftRay := forwardRay.PerpendicularCounterClockwise()
			dright := math.Vector{X: forwardRay.X*math.Cos(45) - math.Sin(45)*forwardRay.Y,
				Y: forwardRay.X*math.Sin(45) + math.Cos(45)*forwardRay.Y}
			dleft := math.Vector{X: forwardRay.X*math.Cos(-45) - math.Sin(-45)*forwardRay.Y,
				Y: forwardRay.X*math.Sin(-45) + math.Cos(-45)*forwardRay.Y}

			world.Buff = world.SpatialHash.Query(math.Square(posComp.Position.Add(forwardRay.MultiplyScalar(movementComp.MovementSpeed*3)), posComp.Radius+150), world.Buff[:0], playerTag, true)

			for _, ent := range world.Buff {
				eIndex := indexMap[ent]
				if eIndex == index {
					continue
				}
				//other := entities[eIndex]
				//me := entities[index]

				otherPos := positionComponents[eIndex].(components.PositionComponent)
				otherMovComp := movementComponents[eIndex].(components.MovementComponent)

				if otherMovComp.Velocity != math.Zero() || entities[eIndex].State != constants.StateAttacking  {
					continue
				}

				detectionLimit := posComp.Radius + otherPos.Radius + 64

				distance := math.GetDistance(posComp.Position.Add(forwardRay.MultiplyScalar(movementComp.MovementSpeed*4)), otherPos.Position)

				if distance < detectionLimit { // cant keep going forard will collide
					ldiagPos := posComp.Position.Add(dleft.MultiplyScalar(movementComp.MovementSpeed*4))
					rdiagPos := posComp.Position.Add(dright.MultiplyScalar(movementComp.MovementSpeed*4))

					distLeft := math.GetDistance(ldiagPos, otherPos.Position)
					distRight := math.GetDistance(rdiagPos, otherPos.Position)

					if distLeft < detectionLimit { // will collide if we turn left
						if distRight < detectionLimit { // cant turn right either, check sides
							velocity = velocity.MultiplyScalar(-1)
							movementComp.DestinationMultiplier = 0.0
						} else { // turn right
							velocity = velocity.Add(posComp.Position.To(rdiagPos).Normalize().MultiplyScalar(3))
							continue
						}
					} else { // can turn left
						if distRight < detectionLimit { // cant turn right, so turn left
							velocity = velocity.Add(posComp.Position.To(ldiagPos).Normalize().MultiplyScalar(3))
							continue
						} else { // check if right is closer

						}
					}
				}

			}
		}

		if movementComp.DestinationMultiplier != 0.0 {
			velocity = velocity.Add(targetDir.MultiplyScalar(movementComp.DestinationMultiplier))
		}

		movementComp.DestinationMultiplier += 0.2

		//Arriving
		distanceToTarget := posComp.Position.Distance(targetPos)
		arrivingZone := posComp.Radius + 80
		//
		//if distanceToTarget < arrivingZone {
		//	desiredVelocity =
		//}

		diff := distanceToTarget / arrivingZone
		if diff < 0.4 * 1/movementComp.MovementSpeed{
			diff = 0.4 * 1/movementComp.MovementSpeed
		} else if diff > 1.0 {
			diff = 1.0
		}
		velocity = velocity.Normalize().MultiplyScalar(movementComp.MovementSpeed * diff)

	//	fmt.Printf("Velocity for %s is %v\n", entities[index].UnitID, velocity)

		if !checkIfUnitInsideMap(posComp.Position.Add(velocity), posComp.Radius/2 - 2) {
			velocity = math.Zero()
			entities[index].State = constants.StateThinking
		}

		posComp.Position = posComp.Position.Add(velocity)

		data := make(map[string]interface{}, 3)
		data["event"] = "walk"
		data["who"] = entities[index].ID
		data["where"] = posComp.Position
		world.ClientEventManager.AddEvent(data)

		movementComp.Velocity = velocity

		posComp.Address = world.SpatialHash.Update(posComp.Address, posComp.Position, entities[index].ID, playerTag)

		world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movementComp
	}

	//fmt.Println("we made ", useless, " iterations in last second")
}

func checkIfUnitInsideMap(pos math.Vector, radius float32) bool {
	isInsideMap := pos.X+radius < float32(constants.MapWidth) && pos.X-radius >= 0 && pos.Y+radius < float32(constants.MapHeight) &&
		pos.Y-radius >= 0

	return isInsideMap
}

