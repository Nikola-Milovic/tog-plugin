package systems

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
	"sort"
)

type Interval struct {
	Start float32
	End   float32
}

type IntervalList []Interval

type MovementSystem struct {
	World       *game.World
	SpatialBuff IntervalList
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

		spatialCheck := true
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

		///	oppositeVec := targetPos.To(posComp.Position)
		//orto := toTarget.Normalize().PerpendicularClockwise()
		//	headingBlocked := false
		if spatialCheck {
			desiredVel := checkAhead(ms.SpatialBuff, entities[index].ID, world)
			if desiredVel != math.Zero() {
				movementComp.DesiredVel = desiredVel.MultiplyScalar(20)
				desiredVel = desiredVel.MultiplyScalar(movementComp.MovementSpeed).Subtract(velocity)
				desiredVel = desiredVel.Normalize().MultiplyScalar(1.2)
				acceleration = acceleration.Add(desiredVel)
				//	velocity = desiredVel
			} else {
				velocity = toTarget.Normalize().MultiplyScalar(movementComp.GoalMultiplier)
			}
		}

		acceleration = acceleration.Add(toTarget.Normalize().MultiplyScalar(0.1))


		diff := distanceToTarget / arrivingZone
		if diff < 0.2*1/movementComp.MovementSpeed {
			diff = 0.2 * 1 / movementComp.MovementSpeed
		} else if diff > 1.0 {
			diff = 1.0
		}

		velocity = velocity.Add(acceleration)
		velocity = velocity.Normalize().MultiplyScalar(movementComp.MovementSpeed * diff)

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

func checkAhead(buff IntervalList, entID int, world *game.World) math.Vector {
	buff = projectAllObstacles(buff, entID, world)
	//	fmt.Printf("Buff is %v\n", buff)
	merged := merge(buff)
	//	sort.Sort(buff)

	fmt.Printf("Merged is %v\n\n", merged)
	for _, interval := range merged {
		if interval.Start < 0 && interval.End > 0 {
			a := float32(0)

			if math.Abs(interval.Start) < math.Abs(interval.End) {
				a = interval.Start
			} else {
				a = interval.End
			}

			vec := math.Vector{X: math.Cos(a) - math.Sin(a),
				Y: math.Sin(a) + math.Cos(a)}
			return vec
		}
	}

	return math.Zero()
}

func projectAllObstacles(buff IntervalList, entID int, world *game.World) IntervalList {
	indexMap := world.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	movComp := movementComponents[indexMap[entID]].(components.MovementComponent)
	posComp := positionComponents[indexMap[entID]].(components.PositionComponent)
	ent := entities[indexMap[entID]]

	targetPos := movComp.Goal
	toTarget := posComp.Position.To(targetPos)
	distanceToTarget := posComp.Position.Distance(targetPos)

	oppositeVec := targetPos.To(posComp.Position)
	//orto := toTarget.Normalize().PerpendicularClockwise()
	defaultAngle := oppositeVec.AngleTo(toTarget) * 57.2957795

	square := getSpatialSquare(posComp.Position, toTarget, posComp.Radius, distanceToTarget)
	world.Buff = world.SpatialHash.Query(square, world.Buff[:0], ent.PlayerTag)

	for _, id := range world.Buff {
		if id == entID {
			continue
		}
		otherIndex := indexMap[id]
		otherEntity := entities[otherIndex]
		otherPosComp := positionComponents[otherIndex].(components.PositionComponent)
		otherMovComp := movementComponents[otherIndex].(components.MovementComponent)

		modifier := float32(1)

		if distanceToTarget < posComp.Position.Distance(otherPosComp.Position) {
			continue
		}

		if isMoving(otherMovComp.Velocity) {
			if otherEntity.State != constants.StateEngaging {
				if otherPosComp.Radius == posComp.Radius { // we are equal, check speed
					if math.Abs(otherMovComp.Velocity.AngleTo(toTarget)) < 5 {
						continue
					}
					if otherMovComp.MovementSpeed > movComp.MovementSpeed { // they are faster, we should just slow down a bit
						modifier = 0.2
					} else if otherMovComp.MovementSpeed < movComp.MovementSpeed { // we are faster they should move
						modifier = 0.1
					} else { // we are equal, just seperate a bit

					}
				} else if otherPosComp.Radius < posComp.Radius { // we are bigger, they move
					continue
				} else { // they are bigger we go around

				}
			}
		} else {
			modifier = 0.6
		}

		tanA, tanB, found := GetTangents(otherPosComp.Position, modifier*(otherPosComp.Radius+posComp.Radius+4), posComp.Position)
		if found {

			tanA = toLocal(tanA, posComp.Position)
			tanB = toLocal(tanB, posComp.Position)

			angle1 := oppositeVec.AngleTo(tanA)*57.2957795 - defaultAngle
			angle2 := oppositeVec.AngleTo(tanB)*57.2957795 - defaultAngle

			interval := Interval{}
			if angle1 > angle2 {
				interval.Start = angle2
				interval.End = angle1
			} else {
				interval.Start = angle1
				interval.End = angle2
			}
			buff = append(buff, interval)
		}
	}

	return buff
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
	return math.Square(position.Add(dirToTarget.Normalize().MultiplyScalar(r/3)), r+2)
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

//Len, Less, Swap needed to execute sort in go
func (list IntervalList) Len() int {
	return len(list)
}
func (list IntervalList) Less(a, b int) bool {
	return list[a].Start < list[b].Start
}

func (list IntervalList) Swap(a, b int) {
	list[a], list[b] = list[b], list[a]
}

func merge(intervals []Interval) []Interval {
	// if empty list
	if len(intervals) <= 2 {
		return intervals
	}
	// As discussed in interview sorting the list first is only way for a O(n) solution below
	sort.Sort(IntervalList(intervals))

	merged := make([]Interval, 0)
	a := &intervals[0]
	listLength := len(intervals)

	for i := 1; i < listLength; i++ {
		b := &intervals[i]
		if a.End >= b.Start {
			if a.End < b.End {
				a.End = b.End
			}
		} else {
			merged = append(merged, *a)
			a = b
		}
	}

	merged = append(merged, *a)
	return merged
}
