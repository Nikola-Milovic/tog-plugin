package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type MovementSystem struct {
	World       *game.World
	SpatialBuff []float32
}

func (ms MovementSystem) Update() {
	world := ms.World
	//useless := 0
	//g := World.Grid
	indexMap := world.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	if ms.SpatialBuff == nil {
		ms.SpatialBuff = make([]float32, 0, 16)
	}

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		if entities[index].State != constants.StateWalking && entities[index].State != constants.StateEngaging {
			movementComp := movementComponents[index].(components.MovementComponent)
			movementComp.Velocity = math.Zero()
			world.ObjectPool.Components["MovementComponent"][index] = movementComp
			continue
		}

		spatialCheck := false
		if world.Tick%4 == 0 {
			spatialCheck = true
		}

		ms.SpatialBuff = ms.SpatialBuff[:0]

		//Setup
		playerTag := entities[index].PlayerTag

		movementComp := movementComponents[index].(components.MovementComponent)
		posComp := positionComponents[index].(components.PositionComponent)
		velocity := movementComp.Velocity
		acceleration := movementComp.Acceleration

		targetPos := movementComp.Goal
		toTarget := posComp.Position.To(targetPos)
		distanceToTarget := posComp.Position.Distance(targetPos)
		//Arriving
		arrivingZone := posComp.Radius + 80


		oppositeVec := targetPos.To(posComp.Position)
		//orto := toTarget.Normalize().PerpendicularClockwise()
		defaultAngle := oppositeVec.AngleTo(toTarget) * 57.2957795
		headingBlocked := false
		square := getSpatialSquare(posComp.Position, toTarget, posComp.Radius, distanceToTarget)
		world.Buff = world.SpatialHash.Query(square, world.Buff[:0], playerTag, true)

		if spatialCheck{
			for _, id := range world.Buff {
				if id == entities[index].ID {
					continue
				}
				otherIndex := indexMap[id]
				//otherEntity := entities[otherIndex]
				otherPosComp := positionComponents[otherIndex].(components.PositionComponent)
				otherMovComp := movementComponents[otherIndex].(components.MovementComponent)

				if isMoving(otherMovComp.Velocity) {
					if math.Abs(otherMovComp.Velocity.AngleTo(toTarget)) < 10 {
						continue
					}
				}

				tanA, tanB, found := GetTangents(otherPosComp.Position, otherPosComp.Radius+posComp.Radius+4, posComp.Position)
				if found {

					tanA = toLocal(tanA, posComp.Position)
					tanB = toLocal(tanB, posComp.Position)

					angle1 := oppositeVec.AngleTo(tanA)*57.2957795 - defaultAngle
					angle2 := oppositeVec.AngleTo(tanB)*57.2957795 - defaultAngle

					if angle1 == 0 || angle2 == 0 || ((angle1 < 0 && angle2 > 0) || (angle1 > 0 && angle2 < 0)) {
						headingBlocked = true
					}

					ms.SpatialBuff = append(ms.SpatialBuff, angle1, angle2)
				}

			}

			minL := float32(180)
			minR := float32(180)
			if headingBlocked {
				for _, val := range ms.SpatialBuff {
					if val < 0 {
						v := val * -1
						if v < minL {
							minL = v
						}
					} else {
						if val < minR {
							minR = val
						}
					}
				}
			}

			//x2=cosβx1−sinβy1
			//y2=sinβx1+cosβy1

			if headingBlocked {
				velocityAngle := oppositeVec.AngleTo(velocity)
				if math.Abs(velocityAngle-(minL*-1)) < math.Abs(velocityAngle-minR) {
					a := minL * -1
					acceleration = acceleration.Add(math.Vector{X: math.Cos(a) - math.Sin(a),
						Y: math.Sin(a) + math.Cos(a)}.MultiplyScalar(1.5))
				} else {
					a := minR
					acceleration = acceleration.Add(math.Vector{X: math.Cos(a) - math.Sin(a),
						Y: math.Sin(a) + math.Cos(a)}.MultiplyScalar(1.5))
				}
			}
		} else {
			acceleration = acceleration.Add(toTarget.Normalize().MultiplyScalar(0.1))
		}


		diff := distanceToTarget / arrivingZone
		if diff < 0.4*1/movementComp.MovementSpeed {
			diff = 0.4 * 1 / movementComp.MovementSpeed
		} else if diff > 1.0 {
			diff = 1.0
		}

		acceleration = acceleration.Add(toTarget.Normalize().MultiplyScalar(0.3))
		velocity = velocity.Add(acceleration)
		velocity = velocity.Normalize().MultiplyScalar(movementComp.MovementSpeed * diff)

		//	fmt.Printf("Velocity for %d is %v\n", entities[index].ID, velocity)
		if !checkIfUnitInsideMap(posComp.Position.Add(velocity), posComp.Radius/2-2) {
			velocity = math.Zero()
			entities[index].State = constants.StateThinking
		}

		//Final
		posComp.Position = posComp.Position.Add(velocity)

		data := make(map[string]interface{}, 3)
		data["event"] = "walk"
		data["who"] = entities[index].ID
		data["where"] = posComp.Position
		world.ClientEventManager.AddEvent(data)

		movementComp.Velocity = velocity
		movementComp.Acceleration = math.Zero()

		posComp.Address = world.SpatialHash.Update(posComp.Address, posComp.Position, entities[index].ID, playerTag)

		world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movementComp
	}

	//fmt.Println("we made ", useless, " iterations in last second")
}

func toLocal(a math.Vector, newZero math.Vector) math.Vector {
	return math.Vector{X: a.X - newZero.X, Y: a.Y - newZero.Y}
}

func checkIfUnitInsideMap(pos math.Vector, radius float32) bool {
	isInsideMap := pos.X+radius < float32(constants.MapWidth) && pos.X-radius >= 0 && pos.Y+radius < float32(constants.MapHeight) &&
		pos.Y-radius >= 0

	return isInsideMap
}

func getSpatialSquare(position, dirToTarget math.Vector, radius, distToTarget float32) math.AABB {
	r := radius + 60
//	center := r/2
	if distToTarget < r {
		r = distToTarget + 2
	//	center = r/2
	}
	return math.Square(position.Add(dirToTarget.Normalize().MultiplyScalar(r/3)), r)
}

func GetSpatalSquareDebug(position, dirToTarget math.Vector, radius, distToTarget float32) math.AABB {
	return getSpatialSquare(position, dirToTarget, radius, distToTarget)
}

//https://github.com/elenzil/tangents/blob/master/tangents/Assets/Scripts/TangentCtlr.cs
//https://answers.unity.com/questions/1617078/finding-a-tangent-vector-from-a-given-point-and-ci.html
func GetTangents(center math.Vector, r float32, p math.Vector) (math.Vector, math.Vector, bool) {
	tanA := math.Zero()
	tanB := math.Zero()

	p = center.To(p)

	dist := p.Magnitute()

	if dist <= r {
		return math.Zero(), math.Zero(), false
	}

	a := r * r / dist
	q := r * math.Sqrt((dist*dist)-(r*r)) / dist

	pN := p.DivideScalar(dist)
	pNP := math.Vector{X: -pN.Y, Y: pN.X}
	va := pN.MultiplyScalar(a)

	tanA = va.Add(pNP.MultiplyScalar(q))
	tanB = va.Subtract(pNP.MultiplyScalar(q))

	tanA = tanA.Add(center)
	tanB = tanB.Add(center)

	return tanA, tanB, true
}

func isMoving(vector math.Vector) bool {
	return vector.X != 0 || vector.Y != 0
}
