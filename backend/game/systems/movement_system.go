package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)


type MovementSystem struct {
	World       *game.World
}

func (ms MovementSystem) Update() {
	world := ms.World
	//useless := 0
	//g := World.Grid
	//indexMap := world.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]
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

		diff := distanceToTarget / arrivingZone
		if diff < 0.2*1/movementComp.MovementSpeed {
			diff = 0.2 * 1 / movementComp.MovementSpeed
		} else if diff > 1.0 {
			diff = 1.0
		}

		desVel := math.Zero()
		desVel = desVel.Add(toTarget.Normalize().MultiplyScalar(0.2))
		desVel = desVel.Add(movementComp.Avoidance)
		desVel = desVel.Add(acceleration)
		desVel = desVel.Normalize().MultiplyScalar(movementComp.MovementSpeed * diff)

		steer := desVel.Subtract(velocity)

		velocity = velocity.Add(steer.Normalize().MultiplyScalar(1.2))

		//	fmt.Printf("Velocity for %d is %v\n", entities[index].ID, velocity)
		if !checkIfUnitInsideMap(posComp.Position.Add(velocity), posComp.Radius/2-2) {
			velocity = posComp.Position.To(math.Vector{X: float32(400), Y: float32(512)}).Normalize()
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
