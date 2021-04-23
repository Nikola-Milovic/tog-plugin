package systems

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/grid"
	"github.com/Nikola-Milovic/tog-plugin/math"
	"strings"
)

type MovementSystem struct {
	World *game.World
	Buff  []string
}

const (
	RepelCof    = 7
	AlignCof    = 0.045
	CohesionCof = 0.03
	Sight       = 100
)

func (ms MovementSystem) Update() {
	world := ms.World
	//useless := 0
	//g := world.Grid
	indexMap := world.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	if ms.Buff == nil {
		ms.Buff = make([]string, 100)
	}

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		movementComp := movementComponents[index].(components.MovementComponent)
		positionComp := positionComponents[index].(components.PositionComponent)

		velocity := movementComp.Velocity

		velocity = velocity.Add(positionComp.Position.To(movementComp.DesiredDestination).Normalize().MultiplyScalar(3))

		//playerTag := entities[index].PlayerTag

		//isMoving := true

		count := float32(1.0)

		// tree rules
		var repel math.Vector
		cohesion := positionComp.Position
		alignment := movementComp.Velocity

		// Now we are querying for ids that are nearby the boid (all tiles that intersect rectangle)
		for _, id := range ms.Buff {
			if entities[index].ID == id {
				continue
			}
			other := &entities[indexMap[id]]
			otherPosComp := positionComponents[other.Index].(components.PositionComponent)
			otherMovementComp := movementComponents[other.Index].(components.MovementComponent)
			dif := otherPosComp.Position.To(positionComp.Position)

			len2 := dif.Len2()

			//if math.Sqrt(len2) > Sight {
			//	//useless++
			//	continue
			//}

			count++
			repel.Add(dif.DivideScalar(len2))
			alignment.Add(otherMovementComp.Velocity)
			cohesion.Add(otherPosComp.Position)
		}

		cohesion = positionComp.Position.To(cohesion.DivideScalar(count)).MultiplyScalar(CohesionCof)
		alignment = alignment.MultiplyScalar(AlignCof / count)

		repel = repel.MultiplyScalar(RepelCof)

		velocity = velocity.Add(cohesion.Add(alignment).Add(repel))
		velocity = limit(velocity, movementComp.MovementSpeed)

		if !checkIfUnitInsideMap(positionComp.Position.Add(velocity), positionComp.BoundingBox) { // so that boids will come back
			velocity = math.Zero()
		}

		positionComp.Position = positionComp.Position.Add(velocity)

		fmt.Printf("I %d am at %v, desired position %v\n", index, positionComp.Position, movementComp.DesiredDestination)

		data := make(map[string]interface{}, 3)
		data["event"] = "walk"
		data["who"] = entities[index].ID
		data["where"] = positionComp.Position
		world.ClientEventManager.AddEvent(data)

		movementComp.Velocity = velocity

		world.ObjectPool.Components["PositionComponent"][index] = positionComp
		world.ObjectPool.Components["MovementComponent"][index] = movementComp
	}

	//fmt.Println("we made ", useless, " iterations in last second")
}

func checkIfUnitInsideMap(pos math.Vector, bbox math.Vector) bool {
	isInsideMap := pos.X+bbox.X < float32(grid.MapWidth) && pos.X-bbox.X >= 0 && pos.Y+bbox.Y < float32(grid.MapHeight) &&
		pos.Y-bbox.Y >= 0

	return isInsideMap
}

func writeImapToConsole(imap *engine.Imap, tick int) {
	var sb strings.Builder

	heading := fmt.Sprintf("Proximity Map Key : %d\n", tick)
	sb.WriteString(heading)
	for y := 0; y < imap.Height; y++ {
		for x := 0; x < imap.Width; x++ {
			s := fmt.Sprintf("%.2f ", imap.Grid[x][y])
			sb.WriteString(s)
		}
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
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
//if world.Tick == 1 || world.Tick%3 == 0 {
//switch entities[index].State {
//
//case constants.StateWalking:
//{
//destination := movementComp.DesiredDestination
//
////	nearbyEntities := helper.GetNearbyEntities(200, world, index)
//
//velocity = velocity.Add(destination.Subtract(positionComp.Position).Normalize().MultiplyScalar(1.3))
//
//avoidance := math.Zero()
////	avoidance = avoidance.Add(alignment(world, nearbyEntities, velocity).MultiplyScalar(alignmentCoef))
//////	avoidance = avoidance.Add(cohesion(world, nearbyEntities, velocity, positionComp.Position).MultiplyScalar(cohesionCoef))
////	avoidance = avoidance.Add(separation(world, nearbyEntities, velocity, positionComp.Position).MultiplyScalar(3))
//
//velocity = velocity.Add(avoidance)
//
//fmt.Printf("I %d am at %v, desired position %v\n\n", index, positionComp.Position, destination)
//}
//case constants.StateEngaging:
//{
//fmt.Printf("Engaging %d \n", index)
////find safe spot around that entity
//velocity = math.Zero()
//}
//default:
//{ //not moving
////isMoving = false
//velocity = math.Zero()
//}
//}
//}
