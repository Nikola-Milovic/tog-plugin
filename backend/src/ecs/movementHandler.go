package ecs

import (
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

// https://github.com/PashaWNN/boids_go/blob/master/boids.go

var alignmentCoef = float32(1.0)
var cohesionCoef = float32(1.0)
var divisionCoef = float32(1.0)
var maxSpeed = float32(0)
var maxForce = float32(0.3)
var desiredSeperation = float32(40)
var breakingForce = float32(0.4)
var avg constants.V2
var total float32

//MovementHandler is a handler used to handle Movement of the entities, Handles the MovementAction
//Calculates the next position an entity should be at
type MovementHandler struct {
	manager *EntityManager
}

//HandleAction handles Movement Action for entity at the given index
func (h MovementHandler) HandleAction(index int) {
	// action, ok := h.manager.Actions[index].(action.MovementAction)

	// if !ok {
	// 	fmt.Println("Error")
	// }

	// destination := action.Target

}
