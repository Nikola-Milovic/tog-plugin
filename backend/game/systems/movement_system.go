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
		ms.SpatialBuff = make([]float32, 0, 20)
	}

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		if entities[index].State != constants.StateWalking && entities[index].State != constants.StateEngaging {
			continue
		}

		ms.SpatialBuff = ms.SpatialBuff[:0]

		//Setup
		playerTag := entities[index].PlayerTag

		movementComp := movementComponents[index].(components.MovementComponent)
		posComp := positionComponents[index].(components.PositionComponent)
		velocity := movementComp.Velocity

		targetPos := movementComp.Goal
		toTarget := posComp.Position.To(targetPos)

		//oppositeVec := targetPos.To(posComp.Position).Normalize()
		orto := toTarget.Normalize().PerpendicularClockwise()

		world.Buff = world.SpatialHash.Query(math.Square(posComp.Position.Add(toTarget.Normalize().MultiplyScalar(64)), posComp.Radius+120), world.Buff[:0], playerTag, true)

		for _, id := range world.Buff {
			if id == entities[index].ID {
				continue
			}
			otherIndex := indexMap[id]
			//otherEntity := entities[otherIndex]
			otherPosComp := positionComponents[otherIndex].(components.PositionComponent)
			//otherMovComp := positionComponents[otherIndex].(components.PositionComponent)

			tanA, tanB, found := GetTangents(otherPosComp.Position, otherPosComp.Radius + 32, posComp.Position)
			if found {
				//fmt.Printf("Angle a is %.2f, angle b is %.2f \n", toTarget.AngleTo(tanA), toTarget.AngleTo(tanB))
				ms.SpatialBuff = append(ms.SpatialBuff, orto.AngleTo(tanA)  * 180/math.Pi - 90,
					orto.AngleTo(tanB) * 180/math.Pi - 90)
			}

		}

		if index == 0 {
			fmt.Println(ms.SpatialBuff)
		}
		if movementComp.DestinationMultiplier != 0.0 {
			velocity = velocity.Add(toTarget.MultiplyScalar(movementComp.DestinationMultiplier))
		}

		movementComp.DestinationMultiplier += 0.2

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
