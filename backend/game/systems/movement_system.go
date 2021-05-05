package systems

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type MovementSystem struct {
	World       *game.World
	SpatialBuff []bool
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
		ms.SpatialBuff = make([]bool, 180, 180)
	}

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		if entities[index].State != constants.StateWalking && entities[index].State != constants.StateEngaging {
			continue
		}

		for indx, _ := range ms.SpatialBuff {
			ms.SpatialBuff[indx] = false
		}

		//Setup
		playerTag := entities[index].PlayerTag

		movementComp := movementComponents[index].(components.MovementComponent)
		posComp := positionComponents[index].(components.PositionComponent)
		velocity := movementComp.Velocity

		targetPos := movementComp.Goal
		toTarget := posComp.Position.To(targetPos)

		//oppositeVec := targetPos.To(posComp.Position)
		orto := toTarget.Normalize().PerpendicularClockwise()
		defaultAngle := orto.AngleTo(toTarget) * 57.2957795

		world.Buff = world.SpatialHash.Query(math.Square(posComp.Position.Add(toTarget.Normalize().MultiplyScalar(64)), posComp.Radius+120), world.Buff[:0], playerTag, true)

		for _, id := range world.Buff {
			if id == entities[index].ID {
				continue
			}
			otherIndex := indexMap[id]
			//otherEntity := entities[otherIndex]
			otherPosComp := positionComponents[otherIndex].(components.PositionComponent)
			//otherMovComp := positionComponents[otherIndex].(components.PositionComponent)

			tanA, tanB, found := GetTangents(otherPosComp.Position, otherPosComp.Radius+16, posComp.Position)
			if found {

				angle1 := orto.AngleTo(tanA)*57.2957795 - defaultAngle
				angle2 := orto.AngleTo(tanB)*57.2957795 - defaultAngle

				//	zeroAngle := orto.AngleTo(toTarget)* 57.2957795 - defaultAngle
				//ms.SpatialBuff = append(ms.SpatialBuff, angle1,
				//	angle2)

				if angle1 < -90 || angle1 > 90 || angle2 < -90 || angle2 > 90 {
					continue
				}

				if angle1 < angle2 {
					for x := int(angle1); x < int(angle2); x++ {
						ms.SpatialBuff[x+89] = true
					}
				} else {
					for x := int(angle2); x < int(angle1); x++ {
						ms.SpatialBuff[x+89] = true
					}
				}

				//if index == 0 {
				//	fmt.Printf("Angle a is %.2f, angle b is %.2f and default is %.2f, zero angle is %.2f\n", angle1,
				//		angle2, defaultAngle,zeroAngle )
				//	fmt.Printf("Tan a is %v, tan b is %vf\n", tanA,
				//		tanB)
				//}

			}

		}

		//if index == 0 {
		//	fmt.Println(ms.SpatialBuff)
		//}


		closestLeft := 0
		closestRight := 179
		for x := 0; x < 90; x++ {
			l := ms.SpatialBuff[x]
			r := ms.SpatialBuff[179-x]

			if !l {
				closestLeft = x
			}
			if !r {
				closestRight = 179 - x
			}
		}

		if velocity == math.Zero() {
			velocity = toTarget.Normalize()
		}

		closestLeft -= 90
		closestRight -= 90

		if math.Absi(89-closestLeft) > math.Absi(89-closestRight) {
			//best angle is closest right
			tarDir := math.Vector{X: velocity.X * math.Cos(float32(closestRight)) - math.Sin(float32(closestRight))*velocity.Y,
				Y: velocity.X * math.Sin(float32(closestRight)) + math.Cos(float32(closestRight))*velocity.Y }
			velocity = velocity.Add(velocity.To(tarDir).Normalize().MultiplyScalar(2))

			fmt.Printf("Closest afree angle is right%d\n", closestRight)
		} else {
			//best angle is closest left
			tarDir := math.Vector{X: velocity.X * math.Cos(float32(closestLeft)) - math.Sin(float32(closestLeft))*velocity.Y,
				Y: velocity.X * math.Sin(float32(closestLeft)) + math.Cos(float32(closestLeft))*velocity.Y }

			velocity = velocity.Add(velocity.To(tarDir).Normalize().MultiplyScalar(2))

			fmt.Printf("Closest afree angle is left %d\n", closestLeft)
		}



		//Arriving
		distanceToTarget := posComp.Position.Distance(targetPos)
		arrivingZone := posComp.Radius + 80

		diff := distanceToTarget / arrivingZone
		if diff < 0.4*1/movementComp.MovementSpeed {
			diff = 0.4 * 1 / movementComp.MovementSpeed
		} else if diff > 1.0 {
			diff = 1.0
		}

		velocity = velocity.Normalize().MultiplyScalar(movementComp.MovementSpeed * diff)

		//	fmt.Printf("Velocity for %s is %v\n", entities[index].UnitID, velocity)

		if !checkIfUnitInsideMap(posComp.Position.Add(velocity), posComp.Radius/2-2) {
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
