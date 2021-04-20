package systems

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
)

// https://github.com/PashaWNN/boids_go/blob/master/boids.go

var alignmentCoef = float32(1.0)
var cohesionCoef = float32(1.0)
var separationCoef = float32(1.5)
var maxSpeed = float32(0)
var maxForce = float32(0.3)
var desiredSeperation = float32(60)
var breakingForce = float32(0.4)

type MovementSystem struct {
	World *game.World
}

func (ms MovementSystem) Update() {
	world := ms.World
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	for index, ent := range entities {
		movementComp := movementComponents[index].(components.MovementComponent)
		positionComp := positionComponents[index].(components.PositionComponent)

		if !movementComp.IsMoving || !entities[index].Active {
			movementComp.Velocity = engine.Zero()
			world.ObjectPool.Components["MovementComponent"][index] = movementComp
			continue
		}

		//Setup
		desiredDirection := movementComp.DesiredDirection //world.Grid.GetDesiredDirectionAt(positionComp.Position, tag)
		velocity := movementComp.Velocity
		maxSpeed = movementComp.MovementSpeed
		desiredSeperation = positionComp.BoundingBox.X/2 + 10

		velocity = velocity.Add(desiredDirection)

		//Avoidance
		nearbyEntities := helper.GetTaggedNearbyEntities(150, world, index, ent.PlayerTag)
		avoidance := engine.Zero()
		avoidance = avoidance.Add(alignment(world, nearbyEntities, velocity))
		avoidance = avoidance.Add(separation(world, nearbyEntities, velocity, positionComp.Position))

		velocity = velocity.Add(avoidance)
		velocity = limit(velocity, maxSpeed)
		positionComp.Position = positionComp.Position.Add(velocity)

		//Finish everything and store the data
		positionComp.Position.X = engine.Constraint(positionComp.Position.X, 0, 799)
		positionComp.Position.Y = engine.Constraint(positionComp.Position.Y, 0, 511)
		movementComp.Velocity = velocity
		world.ObjectPool.Components["PositionComponent"][index] = positionComp
		world.ObjectPool.Components["MovementComponent"][index] = movementComp

		data := make(map[string]interface{}, 3)
		data["event"] = "walk"
		data["who"] = entities[index].ID
		data["where"] = positionComp.Position
		data["velocity"] = movementComp.Velocity
		world.ClientEventManager.AddEvent(data)

		fmt.Printf("I %d am at %v\n", index, positionComp.Position)
	}

}

func limit(p engine.Vector, lim float32) engine.Vector {
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
func alignment(world *game.World, siblings []int, velocity engine.Vector) engine.Vector {
	avg := engine.Vector{X: 0, Y: 0}
	total := float32(0.0)

	for _, siblingIndex := range siblings {
		avg = avg.Add(world.ObjectPool.Components["MovementComponent"][siblingIndex].(components.MovementComponent).Velocity)
		total++
	}
	if total > 0 {
		avg = avg.MultiplyScalar(1 / total * alignmentCoef)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Subtract(velocity)
		avg = limit(avg, maxForce)
		return avg
	}
	return engine.Zero()
}

func cohesion(world *game.World, siblings []int, velocity engine.Vector, position engine.Vector) engine.Vector {
	avg := engine.Vector{X: 0, Y: 0}
	total := float32(0)

	for _, siblingIndex := range siblings {
		avg = avg.Add(world.ObjectPool.Components["PositionComponent"][siblingIndex].(components.PositionComponent).Position)
		total++

	}
	if total > 0 {
		avg = avg.MultiplyScalar(1.0 / total * cohesionCoef)
		avg = avg.Subtract(position)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = limit(avg, maxForce)
		return avg
	}
	return engine.Zero()
}

func separation(world *game.World, siblings []int, velocity engine.Vector, position engine.Vector) engine.Vector {
	avg := engine.Vector{X: 0, Y: 0}
	total := float32(0)

	for _, siblingIndex := range siblings {
		siblingPosComp := world.ObjectPool.Components["PositionComponent"][siblingIndex].(components.PositionComponent)
		siblingPos := siblingPosComp.Position
		d := position.Distance(siblingPos)
		if d < desiredSeperation+siblingPosComp.BoundingBox.X/2 {
			diff := position.Subtract(siblingPos)
			diff = diff.Normalize()
			diff = diff.DivideScalar(d)
			avg = avg.Add(diff)
			total++
		}
	}
	if total > 0 {
		avg = avg.MultiplyScalar(1.0 / total * separationCoef)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Subtract(velocity)
		avg = limit(avg, maxForce)
		return avg
	}
	return engine.Zero()
}
