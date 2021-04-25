package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type MovementSystem struct {
	World *game.World
	Buff  []int
}

var alignmentCoef = float32(1.2)
var cohesionCoef = float32(1.0)
var separationCoef = float32(1.5)
var maxSpeed = float32(0)
var maxForce = float32(0.3)
var desiredSeperation = float32(60)

func (ms MovementSystem) Update() {
	world := ms.World
	//useless := 0
	//g := world.Grid
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

		targetIndex := indexMap[movementComp.TargetID]
		targetPos := positionComponents[targetIndex].(components.PositionComponent)

		desiredSeperation = posComp.Radius

		// Getting all nearby allies
		ms.Buff = world.SpatialHash.Query(math.Square(posComp.Position, 120), ms.Buff[:0], playerTag, true)

		targetPosMultiplier := float32(1.3)

		//Flocking
		avoidance := math.Zero()
		alignmentTotal := float32(0)
		seperationTotal := float32(0)

		avgAlign := math.Zero()
		avgSep := math.Zero()

		for _, entID := range ms.Buff {
			eIndex := indexMap[entID]

			if eIndex == index {
				continue
			}

			entity := entities[eIndex]
			//me := entities[index]

			otherPos := positionComponents[eIndex].(components.PositionComponent)
			otherMov := movementComponents[eIndex].(components.MovementComponent)

			detectionLimit := posComp.Radius + otherPos.Radius
			distance := math.GetDistance(posComp.Position.Add(movementComp.Velocity), otherPos.Position)
			futureDist := math.GetDistance(posComp.Position.Add(movementComp.Velocity.MultiplyScalar(2.5)), otherPos.Position)

			if distance < detectionLimit { // already colliding
				if entity.State == constants.StateAttacking {
					sliding := math.Zero()
					sliding.X = posComp.Position.X - otherPos.Position.X
					sliding.Y = posComp.Position.Y - otherPos.Position.Y

					if sliding.X == 0 {
						sliding.X = 5
					}
					if sliding.Y == 0 {
						sliding.Y = -5
					}

					targetPosMultiplier = 0.0
					velocity = sliding.Normalize().MultiplyScalar(movementComp.MovementSpeed * 0.4)
					break
				} else {
					velocity = math.Zero()
					break
				}
			} else if futureDist < detectionLimit {
				avoid := math.Zero()
				avoid.X = posComp.Position.X - otherPos.Position.X
				avoid.Y = posComp.Position.Y - otherPos.Position.Y

				targetPosMultiplier = 0.2
				velocity = avoidance.Normalize().MultiplyScalar(movementComp.MovementSpeed)
				break
			}

			//Flocking
			avgAlign = avgAlign.Add(otherMov.Velocity)
			if distance < desiredSeperation+otherPos.Radius {
				diff := posComp.Position.Subtract(otherPos.Position)
				diff = diff.Normalize()
				diff = diff.DivideScalar(distance)
				avgSep = avgSep.Add(diff)
				seperationTotal += 1
			}

			alignmentTotal += 1
		}

		avoidance = avoidance.Add(calculateAlignment(alignmentTotal, avgAlign, velocity))
		avoidance = avoidance.Add(calculateSeperation(alignmentTotal, avgAlign, velocity))
		velocity = velocity.Add(avoidance)
		velocity = velocity.Add(posComp.Position.To(targetPos.Position).Normalize().MultiplyScalar(targetPosMultiplier))

		velocity = limit(velocity, movementComp.MovementSpeed)
		if !checkIfUnitInsideMap(posComp.Position.Add(velocity), posComp.Radius) { // so that boids will come back
			velocity = math.Zero()
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

func calculateSeperation(total float32, avg math.Vector, velocity math.Vector) math.Vector {
	if total > 0 {
		avg = avg.DivideScalar(total)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Subtract(velocity)
		avg = limit(avg, 1)
		return avg
	}
	return math.Vector{X: 0.0, Y: 0.0}
}

func calculateAlignment(total float32, avg math.Vector, velocity math.Vector) math.Vector {
	if total > 0 {
		avg = avg.MultiplyScalar(1.0 / total * separationCoef)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Subtract(velocity)
		avg = limit(avg, maxForce)
		return avg
	}
	return math.Zero()
}

func checkIfUnitInsideMap(pos math.Vector, radius float32) bool {
	isInsideMap := pos.X+radius < float32(constants.MapWidth) && pos.X-radius >= 0 && pos.Y+radius < float32(constants.MapHeight) &&
		pos.Y-radius >= 0

	return isInsideMap
}

func limit(p math.Vector, lim float32) math.Vector {
	if p.X > lim {
		p.X = lim
	} else if p.X < -lim {
		p.X = -lim
	}
	if p.Y > lim {
		p.Y = lim
	} else if p.Y < -lim {
		p.Y = -lim
	}
	return p
}

//
//func alignment(world *game.World, siblings []int, velocity math.Vector, id int) math.Vector {
//	avg := math.Vector{X: 0, Y: 0}
//	total := float32(0.0)
//
//	indexMap := world.EntityManager.GetIndexMap()
//
//	for _, siblingId := range siblings {
//		if siblingId == id {
//			continue
//		}
//		avg = avg.Add(world.ObjectPool.Components["MovementComponent"][indexMap[siblingId]].(components.MovementComponent).Velocity)
//		total++
//	}
//	if total > 0 {
//		avg = avg.DivideScalar(total)
//		avg = avg.Normalize().MultiplyScalar(maxSpeed)
//		avg = avg.Subtract(velocity)
//		avg = limit(avg, 1)
//		return avg
//	}
//	return math.Vector{X: 0.0, Y: 0.0}
//
//}
//
//func separation(world *game.World, siblings []int, velocity math.Vector, position math.Vector, id int) math.Vector {
//	avg := math.Vector{X: 0, Y: 0}
//	total := float32(0)
//	indexMap := world.EntityManager.GetIndexMap()
//	for _, siblingId := range siblings {
//		if siblingId == id {
//			continue
//		}
//
//		siblPosComp := world.ObjectPool.Components["PositionComponent"][indexMap[siblingId]].(components.PositionComponent)
//		siblingPos := siblPosComp.Position
//		d := position.Distance(siblingPos)
//		if d < desiredSeperation+siblPosComp.Radius {
//			diff := position.Subtract(siblingPos)
//			diff = diff.Normalize()
//			diff = diff.DivideScalar(d)
//			avg = avg.Add(diff)
//			total++
//		}
//	}
//	if total > 0 {
//		avg.DivideScalar(total)
//	}
//
//	if total > 0 {
//		avg = avg.MultiplyScalar(1.0 / total * separationCoef)
//		avg = avg.Normalize().MultiplyScalar(maxSpeed)
//		avg = avg.Subtract(velocity)
//		avg = limit(avg, maxForce)
//	}
//	return avg
//}
//
