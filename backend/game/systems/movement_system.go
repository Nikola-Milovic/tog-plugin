package systems

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/grid"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
	"github.com/Nikola-Milovic/tog-plugin/math"
	"github.com/Nikola-Milovic/tog-plugin/startup"
	"strings"
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
var breakingForce = float32(0.4)

func (ms MovementSystem) Update() {
	world := ms.World
	//useless := 0
	g := world.Grid
	indexMap := world.GetEntityManager().GetIndexMap()
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

		movementComp := movementComponents[index].(components.MovementComponent)
		positionComp := positionComponents[index].(components.PositionComponent)
		velocity := movementComp.Velocity

		targetIndex := indexMap[movementComp.TargetID]
		targetPos := positionComponents[targetIndex].(components.PositionComponent)

		switch entities[index].State {

		case constants.StateWalking:
			{

				ms.Buff = helper.GetNearbyEntities(150, world, index, ms.Buff[:0])

				velocity = velocity.Add(positionComp.Position.To(targetPos.Position).Normalize().MultiplyScalar(3))

				//isMoving := true

				avoidance := math.Zero()
				avoidance = avoidance.Add(alignment(world, ms.Buff, velocity).MultiplyScalar(alignmentCoef))
				//	avoidance = avoidance.Add(cohesion(world, nearbyEntities, velocity, positionComp.Position).MultiplyScalar(cohesionCoef))
				avoidance = avoidance.Add(separation(world, ms.Buff, velocity, positionComp.Position).MultiplyScalar(2))

				velocity = velocity.Add(avoidance)
			}

		case constants.StateEngaging:
			{

				ms.Buff = helper.GetNearbyEntities(150, world, index, ms.Buff[:0])

				//isMoving := true

				playerTag := entities[index].PlayerTag
				size := int((movementComp.MovementSpeed*50 + 20) / constants.TileSize)
				wm := g.GetWorkingMap(size, size)

				mx, my := grid.GlobalCordToTiled(positionComp.Position.Add(velocity))

				engine.AddIntoSmallerMap(g.GetEnemyProximityImap(playerTag), wm, mx, my, 1)
				//writeImapToConsole(wm, world.Tick)
				engine.AddIntoSmallerMap(g.GetProximityImaps()[playerTag], wm, mx, my, 0.6)
				engine.AddIntoBiggerMap(grid.GetProximityTemplate(movementComp.MovementSpeed).Imap, wm, size/2, size/2, -0.6)

				wm.NormalizeAndInvert()
				engine.AddIntoBiggerMap(startup.InterestTemplates[0].Imap, wm, size/2, size/2, 1)
				//
				x, y, _ := wm.GetLowestValue()

				tarX, tarY := grid.GetBaseMapCoordsFromSectionImapCoords(mx, my, x, y)

				v := positionComp.Position.To(math.Vector{X: float32(tarX * constants.TileSize), Y: float32(tarY * constants.TileSize)}).Normalize().MultiplyScalar(2)

				velocity = velocity.Add(v)

				writeImapToConsole(wm, world.Tick)

				//	fmt.Printf("I %d am at %v, desired position %v\n", index, positionComp.Position, movementComp.TargetID)

				//fromMeToTarget := positionComp.Position.To(targetPos.Position)

				//	velocity = velocity.Add(fromMeToTarget.MultiplyScalar(1.4))

				avoidance := math.Zero()
			//	avoidance = avoidance.Add(alignment(world, ms.Buff, velocity).MultiplyScalar(alignmentCoef))
				//	avoidance = avoidance.Add(cohesion(world, nearbyEntities, velocity, positionComp.Position).MultiplyScalar(cohesionCoef))
				avoidance = avoidance.Add(separation(world, ms.Buff, velocity, positionComp.Position).MultiplyScalar(2))

				velocity = velocity.Add(avoidance)

				//isMoving := true
			}
		default:
			velocity = math.Zero()

		}

		velocity = limit(velocity, movementComp.MovementSpeed)
		if !checkIfUnitInsideMap(positionComp.Position.Add(velocity), positionComp.BoundingBox) { // so that boids will come back
			velocity = math.Zero()
		}

		positionComp.Position = positionComp.Position.Add(velocity)

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

func alignment(world *game.World, siblings []int, velocity math.Vector) math.Vector {
	avg := math.Vector{X: 0, Y: 0}
	total := float32(0.0)

	for _, siblingIndex := range siblings {
		avg = avg.Add(world.ObjectPool.Components["MovementComponent"][siblingIndex].(components.MovementComponent).Velocity)
		total++
	}
	if total > 0 {
		avg = avg.DivideScalar(total)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Subtract(velocity)
		avg = limit(avg, 1)
		return avg
	}
	return math.Vector{X: 0.0, Y: 0.0}

}

func cohesion(world *game.World, siblings []int, velocity math.Vector, position math.Vector) math.Vector {
	avg := math.Vector{X: 0, Y: 0}
	total := float32(0)

	for _, siblingIndex := range siblings {

		avg = avg.Add(world.ObjectPool.Components["PositionComponent"][siblingIndex].(components.PositionComponent).Position)
		total++

	}
	if total > 0 {
		avg = avg.MultiplyScalar(1.0 / total * cohesionCoef)
		avg = avg.Subtract(position)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Subtract(velocity)
		avg = limit(avg, maxForce)
		return avg
	}
	return math.Vector{X: 0.0, Y: 0.0}
}

func separation(world *game.World, siblings []int, velocity math.Vector, position math.Vector) math.Vector {
	avg := math.Vector{X: 0, Y: 0}
	total := float32(0)

	for _, siblingIndex := range siblings {
		siblingPos := world.ObjectPool.Components["PositionComponent"][siblingIndex].(components.PositionComponent).Position
		d := position.Distance(siblingPos)
		if d < desiredSeperation {
			diff := position.Subtract(siblingPos)
			diff = diff.Normalize()
			diff = diff.DivideScalar(d)
			avg = avg.Add(diff)
			total++
		}
	}
	if total > 0 {
		avg.DivideScalar(total)
	}

	if total > 0 {
		avg = avg.MultiplyScalar(1.0 / total * separationCoef)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Subtract(velocity)
		avg = limit(avg, maxForce)
	}
	return avg
}
