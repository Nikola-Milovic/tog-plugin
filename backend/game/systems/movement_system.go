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
	for index, comp := range world.ObjectPool.Components["MovementComponent"] {
		movementComp := comp.(components.MovementComponent)
		positionComp := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
		tag := entities[index].PlayerTag

		if !movementComp.IsMoving || !world.EntityManager.GetEntities()[index].Active {
			movementComp.Velocity = engine.Zero()
			world.ObjectPool.Components["MovementComponent"][index] = movementComp
			continue
		}

		desiredDirection :=  movementComp.DesiredDirection  //world.Grid.GetDesiredDirectionAt(positionComp.Position, tag)
		velocity := movementComp.Velocity
		maxSpeed = movementComp.MovementSpeed
		desiredSeperation = positionComp.BoundingBox.X/2 + 10

		velocity = velocity.Add(desiredDirection)

		//Avoidance
		nearbyEntities := helper.GetTaggedNearbyEntities(100, world, index, tag)


		avoidance := engine.Zero()
		avoidance = avoidance.Add(alignment(world, nearbyEntities, velocity))
	//	avoidance = avoidance.Add(cohesion(world, nearbyEntities, velocity, positionComp.Position))
		avoidance = avoidance.Add(separation(world, nearbyEntities, velocity, positionComp.Position))

		//lookAheadVectorLong := velocity.Add(desiredDirection).MultiplyScalar(maxSpeed * 2.5)
		//lookAheadVectorShort := velocity.Add(desiredDirection).MultiplyScalar(maxSpeed)
		//maxAvoidanceForce := float32(1.0)
		//
		//checkPosShort := positionComp.Position.Add(lookAheadVectorShort)
		//checkPosLong := positionComp.Position.Add(lookAheadVectorLong)
		//
		//collidedIndexShort := world.Grid.IsPositionFree(index, checkPosShort, positionComp.BoundingBox)
		//collidedIndexLong := world.Grid.IsPositionFree(index, checkPosLong, positionComp.BoundingBox)
		//
		//if collidedIndexShort != -1 {
		//	velocity = engine.Zero()
		//	avoidance = checkPosShort.Subtract(world.ObjectPool.Components["PositionComponent"][collidedIndexShort].(components.PositionComponent).Position).Normalize()
		//	avoidance = avoidance.MultiplyScalar(maxAvoidanceForce * 1.5)
		//	//checkPosShort = world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent).Position.Add(lookAheadVectorShort)
		//	//collidedIndexShort = world.Grid.IsPositionFree(index, checkPosShort, positionComp.BoundingBox)
		//} else if collidedIndexLong != -1 {
		//	velocity = velocity.MultiplyScalar(breakingForce)
		//	avoidance = checkPosShort.Subtract(world.ObjectPool.Components["PositionComponent"][collidedIndexLong].(components.PositionComponent).Position).Normalize()
		//	avoidance = avoidance.MultiplyScalar(maxAvoidanceForce * 1.2)
		//	//checkPosLong = world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent).Position.Add(lookAheadVectorLong)
		//	//collidedIndexLong = world.Grid.IsPositionFree(index, checkPosLong, positionComp.BoundingBox)
		//}

		velocity = velocity.Add(avoidance)

		velocity = limit(velocity, maxSpeed)
		positionComp.Position = positionComp.Position.Add(velocity)

		positionComp.Position.X = engine.Constraint(positionComp.Position.X, 0, 799)
		positionComp.Position.Y = engine.Constraint(positionComp.Position.Y, 0, 511)

		movementComp.Velocity = velocity

		world.ObjectPool.Components["PositionComponent"][index] = positionComp
		world.ObjectPool.Components["MovementComponent"][index] = movementComp

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
		avg = avg.MultiplyScalar(1/ total*alignmentCoef)
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
		if d < desiredSeperation + siblingPosComp.BoundingBox.X/2 {
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
