package ecs

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/src/action"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

// https://github.com/PashaWNN/boids_go/blob/master/boids.go

var alignmentCoef = float32(1.0)
var cohesionCoef = float32(1.0)
var divisionCoef = float32(1.0)
var maxSpeed = float32(0)
var maxForce = float32(0.3)
var desiredSeperation = float32(40)
var avg constants.V2
var total float32

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

	nearby := h.manager.getNearbyEntities(50, h.manager.PositionComponents[index].Position, index)

	destination := action.Destination

	direction := destination.Subtract(h.manager.PositionComponents[index].Position).Normalize()

	velocity := direction

	maxSpeed = h.manager.MovementComponents[index].Speed
	// alignmentPerception = float32(h.manager.Entities[index].Size.X / 2)
	// cohesionPerception = float32(h.manager.Entities[index].Size.X/2 + 5)
	// divisionPerception = float32(h.manager.Entities[index].Size.X / 2)

	velocity = velocity.Add(alignment(index, h.manager, nearby, h.manager.MovementComponents[index].Velocity).MultiplyScalar(1.0))
	//velocity = velocity.Add(cohesion(index, h.manager, nearby, h.manager.MovementComponents[index].Velocity).MultiplyScalar(1.0))
	velocity = velocity.Add(separation(index, h.manager, nearby, h.manager.MovementComponents[index].Velocity).MultiplyScalar(1.2))

	

	if(  h.manager.isPositionFree (h.manager.PositionComponents[index].Position.Add(velocity.MultiplyScalar(maxSpeed)))) {
		h.manager.PositionComponents[index].Position = 
	} else {

	}

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

// A method that calculates and applies a steering force towards a target
// STEER = DESIRED MINUS VELOCITY
func seek(target constants.V2, position constants.V2, velocity constants.V2) constants.V2 {
	desired := target.Subtract(position) // A vector pointing from the position to the target
	// Scale to maximum speed
	desired.Normalize()
	desired.MultiplyScalar(maxSpeed)

	// Above two lines of code below could be condensed with new PVector setMag() method
	// Not using this method until Processing.js catches up
	// desired.setMag(maxspeed);

	// Steering = Desired minus Velocity
	steer := desired.Subtract(velocity)
	steer = limit(steer, maxForce) // Limit to maximum steering force
	return steer
}

func alignment(index int, e *EntityManager, siblings []int, velocity constants.V2) constants.V2 {
	avg = constants.V2{X: 0, Y: 0}
	total = 0.0

	for _, sibling := range siblings {
		avg = avg.Add(e.MovementComponents[sibling].Velocity)
		total++
	}
	if total > 0 {
		avg = avg.DivideScalar(total)
		avg = avg.Normalize().MultiplyScalar(maxSpeed)
		avg = avg.Subtract(velocity)
		avg = limit(avg, maxForce)
		return avg
	}
	return constants.New(0, 0)

}

func cohesion(index int, e *EntityManager, siblings []int, velocity constants.V2) constants.V2 {
	avg = constants.V2{X: 0, Y: 0}
	total = 0

	for _, sibling := range siblings {

		avg = avg.Add(e.PositionComponents[sibling].Position)
		total++

	}
	if total > 0 {
		avg = avg.DivideScalar(total)
		return seek(avg, e.PositionComponents[index].Position, velocity)
	}
	return constants.New(0, 0)
}

func separation(index int, e *EntityManager, siblings []int, velocity constants.V2) constants.V2 {
	avg = constants.V2{X: 0, Y: 0}
	total = 0

	for _, sibling := range siblings {
		d := e.PositionComponents[index].Position.Distance(e.PositionComponents[sibling].Position)
		if d < desiredSeperation {
			diff := e.PositionComponents[index].Position.Subtract(e.PositionComponents[sibling].Position)
			diff = diff.Normalize()
			diff = diff.DivideScalar(d)
			avg = avg.Add(diff)
			total++
		}
	}
	if total > 0 {
		avg.DivideScalar(total)
	}

	if avg.Magnitute() > 0 {
		avg.Normalize()
		avg.MultiplyScalar(maxSpeed)
		avg.Subtract(velocity)
		avg = limit(avg, maxForce)
	}
	return avg
}
