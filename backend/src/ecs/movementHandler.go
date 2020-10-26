package ecs

import "github.com/Nikola-Milovic/tog-plugin/src/action"

type MovementHandler struct {
	manager *EntityManager
}

func (h MovementHandler) HandleAction(index int) {
	action, ok := h.manager.Actions[index].(action.MovementAction)

	if !ok {
		println("error")
	}

	destination := action.Destination

	direction := h.manager.PositionComponents[index].Position.Subtract(destination).Normalize()

	h.manager.PositionComponents[index].Position.Add(direction.MultiplyScalar(float64(h.manager.MovementComponents[index].Speed)))
}
