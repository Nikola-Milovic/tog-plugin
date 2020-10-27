package ecs

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/src/action"
)

type MovementHandler struct {
	manager *EntityManager
}

func (h MovementHandler) HandleAction(index int) {
	action, ok := h.manager.Actions[index].(action.MovementAction)

	if !ok {
		fmt.Println("Error")
	}

	destination := action.Destination.MultiplyScalar(float64(index + 1))

	direction := destination.Subtract(h.manager.PositionComponents[index].Position).Normalize()

	h.manager.PositionComponents[index].Position = h.manager.PositionComponents[index].Position.Add((direction.MultiplyScalar(float64(h.manager.MovementComponents[index].Speed))))

	fmt.Println(direction)
	fmt.Println(direction.MultiplyScalar(float64(h.manager.MovementComponents[index].Speed)))
	fmt.Println(h.manager.PositionComponents[index].Position)
}
