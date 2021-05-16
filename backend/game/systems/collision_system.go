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

//CollisionSystem is supposed to do
//1) Not allow overlapping of units
//2) Keep units from bumping into each other
//3) Let some units pass before others/ go around someone
type CollisionSystem struct {
	World       *game.World
	SpatialBuff IntervalList
}

func (ms CollisionSystem) Update() {
	ms.collision()
	ms.pathPlotting()
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

		if movComp.Velocity == math.Zero() {
			continue
		}

		world.Buff = world.SpatialHash.Query(math.Square(posComp.Position, posComp.Radius+80), world.Buff[:0], -1)

		//goal := movComp.Goal

		movComp.GoalMultiplier = 1.0

		//	distToGoal := math.GetDistance(goal, posComp.Position)
		adjustment := math.Zero()

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

				if distance < detectionLimit-4 { // overlapping
					adjustment = otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(1.2)
				}

				if otherMovComp.Velocity != math.Zero() { // seperation while walking
					movComp.GoalMultiplier = 0.7
					if distance < detectionLimit*1.5 {
						adjustment = otherPos.Position.To(posComp.Position).Normalize()
						if otherPos.Radius == posComp.Radius { // we are equal, check speed
							if otherMovComp.MovementSpeed > movComp.MovementSpeed { // they are faster, we should just slow down a bit
								adjustment = adjustment.MultiplyScalar(0.4)
							} else if otherMovComp.MovementSpeed < movComp.MovementSpeed { // we are faster they should move
								adjustment = adjustment.MultiplyScalar(0.3)
							} else { // we are equal, just seperate a bit
								adjustment = adjustment.MultiplyScalar(0.3)
							}
						} else if otherPos.Radius < posComp.Radius { // we are bigger, they move
							adjustment = adjustment.MultiplyScalar(0.2)
						} else { // they are bigger we go around
							adjustment = adjustment.MultiplyScalar(1.2)
						}
					}
				} else {
					if distance < detectionLimit*1.8 { // other is standing
						movComp.GoalMultiplier = 0.2
						adjustment = otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(0.4)
					}
				}


			} else {
				detectionLimit := posComp.Radius + otherPos.Radius + atkComp.Range + 2
				distance := math.GetDistance(posComp.Position, otherPos.Position)

				//Overlap
				if distance < detectionLimit {
					movComp.GoalMultiplier = 0.7
					adjustment = otherPos.Position.To(posComp.Position).Normalize().MultiplyScalar(0.5)
				}

			}
		}

		//if distToGoal < atkComp.Range + posComp.Radius + movComp.MovementSpeed*3 + 16 {
		//	movComp.Acceleration = movComp.Acceleration.MultiplyScalar(0.3)
		//}
		fmt.Printf("Adjustment %v, buff : %d\n", adjustment, len(world.Buff))
		if adjustment.X == 0 && adjustment.Y==0 {
			adjustment =adjustment.MultiplyScalar(0.1)
			movComp.Avoidance = movComp.Avoidance.MultiplyScalar(0.8)
		}
		movComp.Avoidance = movComp.Avoidance.Add(adjustment) //.Add(pushTowardsGoal)
		world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movComp
	}
}

func (ms CollisionSystem) pathPlotting() {
	world := ms.World
	//	indexMap := ms.World.EntityManager.GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	//	attackComponents := World.ObjectPool.Components["AttackComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		//unitTag := entities[index].PlayerTag

		posComp := positionComponents[index].(components.PositionComponent)
		//	atkComp := attackComponents[index].(components.AttackComponent)
		movComp := movementComponents[index].(components.MovementComponent)

		if movComp.Velocity == math.Zero() {
			continue
		}

//		targetPos := movComp.Goal
	//	toTarget := posComp.Position.To(targetPos)

		velocity := movComp.Velocity
		acceleration := movComp.Acceleration

		//desiredVel := checkAhead(ms.SpatialBuff, entities[index].ID, world)
		//oppositeVec := targetPos.To(posComp.Position)
		//defaultAngle := oppositeVec.AngleTo(toTarget) * 57.2957795
		//if desiredVel != math.Zero() {
		//	fmt.Printf("Actual Angle is %.2f, opposite vec is %v\n\n", oppositeVec.AngleTo(desiredVel)*57.2957795-defaultAngle, oppositeVec)
		//	movComp.DesiredVel = desiredVel.MultiplyScalar(20)
		//	desiredVel = desiredVel.MultiplyScalar(movComp.MovementSpeed).Subtract(velocity)
		//	desiredVel = desiredVel.Normalize().MultiplyScalar(0.6)
		//	acceleration = acceleration.Add(desiredVel)
		//	//	velocity = desiredVel
		//} else {
		//	velocity = toTarget.Normalize().MultiplyScalar(movComp.GoalMultiplier)
		//}

		movComp.Velocity = velocity
		movComp.Acceleration = acceleration

		world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movComp
	}
}

func checkAhead(buff IntervalList, entID int, world *game.World) math.Vector {
	buff = projectAllObstacles(buff, entID, world)
	//	fmt.Printf("Buff is %v\n", buff)
	merged := merge(buff)
	//	sort.Sort(buff)

	fmt.Printf("Merged is %v\n", merged)
	for _, interval := range merged {
		if interval.Start < 0 && interval.End > 0 {
			a := float32(0)

			if math.Abs(interval.Start) < math.Abs(interval.End) {
				a = interval.Start
			} else {
				a = interval.End
			}

			fmt.Printf("Desired Angle is %.2f\n", a)
			a = a * math.Pi / 180
			fmt.Printf("After modification is %.2f\n", a)
			vec := math.Vector{X: math.Cos(a),
				Y: math.Sin(a)}

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

	fmt.Printf("Opposite vec is %v \n", oppositeVec)

	square := getSpatialSquare(posComp.Position, toTarget, posComp.Radius, distanceToTarget)
	world.Buff = world.SpatialHash.Query(square, world.Buff[:0], ent.PlayerTag)

	for _, id := range world.Buff {
		if id == entID {
			continue
		}
		otherIndex := indexMap[id]
		//	otherEntity := entities[otherIndex]
		otherPosComp := positionComponents[otherIndex].(components.PositionComponent)
		//	otherMovComp := movementComponents[otherIndex].(components.MovementComponent)

		tanA, tanB, found := GetTangents(otherPosComp.Position, otherPosComp.Radius+posComp.Radius+4, posComp.Position)
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

func getSpatialSquare(position, dirToTarget math.Vector, radius, distToTarget float32) math.AABB {
	r := radius + 60
	//	center := r/2
	if distToTarget < r {
		r = distToTarget + 2
		//	center = r/2
	}
	return math.Square(position.Add(dirToTarget.Normalize().MultiplyScalar(r/3)), r+2)
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
