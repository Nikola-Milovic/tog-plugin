package ecs

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/src/action"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

// https://github.com/PashaWNN/boids_go/blob/master/boids.go

var alignmentPerception = float32(20.0)
var cohesionPerception = float32(40.0)
var divisionPerception = float32(20.0)
var alignmentCoef = float32(1.0)
var cohesionCoef = float32(1.0)
var divisionCoef = float32(1.0)
var maxSpeed = float32(0)
var maxForce = float32(1)
var avg constants.V2
var total int

//MovementHandler is a handler used to handle Movement of the entities, Handles the MovementAction
//Calculates the next position an entity should be at
type MovementHandler struct {
	manager *EntityManager
}

//HandleAction handles Movement Action for entity at the given index
func (h MovementHandler) HandleAction(index int) {
	action, ok := h.manager.Actions[index].(action.MovementAction)

	if !ok {
		fmt.Println("Error")
	}

	destination := action.Destination

	direction := destination.Subtract(h.manager.PositionComponents[index].Position).Normalize()

	maxSpeed = h.manager.MovementComponents[index].Speed
	// alignmentPerception = float32(h.manager.Entities[index].Size.X / 2)
	// cohesionPerception = float32(h.manager.Entities[index].Size.X/2 + 5)
	// divisionPerception = float32(h.manager.Entities[index].Size.X / 2)

	//Alignment
	h.manager.MovementComponents[index].Accelaration = h.manager.MovementComponents[index].Accelaration.Add(alignment(index, h.manager, action.NearbyEntities))

	//Cohesion
	h.manager.MovementComponents[index].Accelaration = h.manager.MovementComponents[index].Accelaration.Add(cohesion(index, h.manager, action.NearbyEntities))

	//Seperation
	h.manager.MovementComponents[index].Accelaration = h.manager.MovementComponents[index].Accelaration.Add(separation(index, h.manager, action.NearbyEntities))

	//Add the target position to our acceleration
	h.manager.MovementComponents[index].Accelaration = h.manager.MovementComponents[index].Accelaration.Add(direction)

	//Add the velocity to the position
	h.manager.PositionComponents[index].Position = h.manager.PositionComponents[index].Position.Add(h.manager.MovementComponents[index].Velocity)
	h.manager.MovementComponents[index].Velocity = h.manager.MovementComponents[index].Velocity.Add(h.manager.MovementComponents[index].Accelaration.MultiplyScalar(0.5))
	h.manager.MovementComponents[index].Velocity = limit(h.manager.MovementComponents[index].Velocity, maxSpeed)
	h.manager.MovementComponents[index].Accelaration = constants.Zero()
}

func limit(p constants.V2, lim float32) constants.V2 {
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

func alignment(index int, e *EntityManager, siblings []int) constants.V2 {
	avg = constants.V2{0, 0}
	total = 0

	for _, sibling := range siblings {
		if e.PositionComponents[index].Position.Distance(e.PositionComponents[sibling].Position) < alignmentPerception {
			avg = avg.Add(e.MovementComponents[sibling].Velocity)
			total++
		}
	}
	if total > 0 {
		avg = avg.MultiplyScalar(1.0 / float32(total) * alignmentCoef)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Divide(e.MovementComponents[index].Velocity)
		avg = limit(avg, maxForce)
	}
	return avg
}

func cohesion(index int, e *EntityManager, siblings []int) constants.V2 {
	avg = constants.V2{0, 0}
	total = 0

	for _, sibling := range siblings {
		if e.PositionComponents[index].Position.Distance(e.PositionComponents[sibling].Position) < cohesionPerception {
			avg = avg.Add(e.PositionComponents[sibling].Position)
			total++
		}
	}
	if total > 0 {
		avg = avg.MultiplyScalar(1.0 / float32(total) * cohesionCoef)
		avg = avg.Divide(e.PositionComponents[index].Position)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Divide(e.MovementComponents[index].Velocity)
		avg = limit(avg, maxForce)
	}
	return avg
}

func separation(index int, e *EntityManager, siblings []int) constants.V2 {
	avg = constants.V2{0, 0}
	total = 0

	for _, sibling := range siblings {
		d := e.PositionComponents[index].Position.Distance(e.PositionComponents[sibling].Position)
		if d < divisionPerception {
			diff := e.PositionComponents[index].Position.Subtract(e.PositionComponents[sibling].Position).MultiplyScalar(1 / (d * d))
			avg = avg.Add(diff)
			total++
		}
	}
	if total > 0 {
		avg = avg.MultiplyScalar(1.0 / float32(total) * divisionCoef)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Subtract(e.MovementComponents[index].Velocity)
		avg = limit(avg, maxForce)
	}
	return avg
}
