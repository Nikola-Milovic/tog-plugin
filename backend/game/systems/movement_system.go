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

		playerTag := entities[index].PlayerTag

		movementComp := movementComponents[index].(components.MovementComponent)
		positionComp := positionComponents[index].(components.PositionComponent)
		velocity := movementComp.Velocity

		targetIndex := indexMap[movementComp.TargetID]
		targetPos := positionComponents[targetIndex].(components.PositionComponent)

		desiredSeperation = positionComp.Radius

		// Now we are querying for ids that are nearby the boid (all tiles that intersect rectangle)
		ms.Buff = world.SpatialHash.Query(math.Square(positionComp.Position, 120), ms.Buff[:0], playerTag, true)

		targetPosMultiplier := float32(1.3)

		velocity = velocity.Add(positionComp.Position.To(targetPos.Position).Normalize().MultiplyScalar(targetPosMultiplier))

		//isMoving := true

		avoidance := math.Zero()
		avoidance = avoidance.Add(alignment(world, ms.Buff, velocity, entities[index].ID).MultiplyScalar(alignmentCoef))
		//	avoidance = avoidance.Add(cohesion(world, nearbyEntities, velocity, positionComp.Position).MultiplyScalar(cohesionCoef))
		avoidance = avoidance.Add(separation(world, ms.Buff, velocity, positionComp.Position, entities[index].ID).MultiplyScalar(2))

		velocity = velocity.Add(avoidance)

		switch entities[index].State {

		case constants.StateWalking:
			{
			}
		case constants.StateEngaging:
			{
			}
		default:
			velocity = math.Zero()
		}

		velocity = limit(velocity, movementComp.MovementSpeed)
		if !checkIfUnitInsideMap(positionComp.Position.Add(velocity), positionComp.Radius) { // so that boids will come back
			velocity = math.Zero()
		}

		positionComp.Position = positionComp.Position.Add(velocity)

		data := make(map[string]interface{}, 3)
		data["event"] = "walk"
		data["who"] = entities[index].ID
		data["where"] = positionComp.Position
		world.ClientEventManager.AddEvent(data)

		movementComp.Velocity = velocity

		positionComp.Address = world.SpatialHash.Update(positionComp.Address, positionComp.Position, entities[index].ID, playerTag)

		world.ObjectPool.Components["PositionComponent"][index] = positionComp
		world.ObjectPool.Components["MovementComponent"][index] = movementComp
	}

	//fmt.Println("we made ", useless, " iterations in last second")
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

func alignment(world *game.World, siblings []int, velocity math.Vector, id int) math.Vector {
	avg := math.Vector{X: 0, Y: 0}
	total := float32(0.0)

	indexMap := world.EntityManager.GetIndexMap()

	for _, siblingId := range siblings {
		if siblingId == id {
			continue
		}
		avg = avg.Add(world.ObjectPool.Components["MovementComponent"][indexMap[siblingId]].(components.MovementComponent).Velocity)
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

func separation(world *game.World, siblings []int, velocity math.Vector, position math.Vector, id int) math.Vector {
	avg := math.Vector{X: 0, Y: 0}
	total := float32(0)
	indexMap := world.EntityManager.GetIndexMap()
	for _, siblingId := range siblings {
		if siblingId == id {
			continue
		}

		siblPosComp := world.ObjectPool.Components["PositionComponent"][indexMap[siblingId]].(components.PositionComponent)
		siblingPos := siblPosComp.Position
		d := position.Distance(siblingPos)
		if d < desiredSeperation+siblPosComp.Radius {
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

//
//playerTag := entities[index].PlayerTag
//size := int((movementComp.MovementSpeed*50 + 20) / constants.TileSize)
//wm := g.GetWorkingMap(size, size)
//
//mx, my := grid.GlobalCordToTiled(positionComp.Position.Add(velocity))
//
//engine.AddIntoSmallerMap(g.GetEnemyProximityImap(playerTag), wm, mx, my, 1)
////writeImapToConsole(wm, world.Tick)
//engine.AddIntoSmallerMap(g.GetProximityImaps()[playerTag], wm, mx, my, 0.6)
//engine.AddIntoBiggerMap(grid.GetProximityTemplate(movementComp.MovementSpeed).Imap, wm, size/2, size/2, -0.6)
//
//wm.NormalizeAndInvert()
//engine.AddIntoBiggerMap(startup.InterestTemplates[0].Imap, wm, size/2, size/2, 1)
////
//x, y, _ := wm.GetLowestValue()
//
//tarX, tarY := grid.GetBaseMapCoordsFromSectionImapCoords(mx, my, x, y)
//
//v := positionComp.Position.To(math.Vector{X: float32(tarX * constants.TileSize), Y: float32(tarY * constants.TileSize)}).Normalize().MultiplyScalar(2)
//
//velocity = velocity.Add(v)
//
//writeImapToConsole(wm, world.Tick)

//func cohesion(world *game.World, siblings []int, velocity math.Vector, position math.Vector) math.Vector {
//	avg := math.Vector{X: 0, Y: 0}
//	total := float32(0)
//
//	for _, siblingIndex := range siblings {
//
//		avg = avg.Add(world.ObjectPool.Components["PositionComponent"][siblingIndex].(components.PositionComponent).Position)
//		total++
//
//	}
//	if total > 0 {
//		avg = avg.MultiplyScalar(1.0 / total * cohesionCoef)
//		avg = avg.Subtract(position)
//		avg = avg.Normalize().MultiplyScalar(maxSpeed)
//		avg = avg.Subtract(velocity)
//		avg = limit(avg, maxForce)
//		return avg
//	}
//	return math.Vector{X: 0.0, Y: 0.0}
//}
