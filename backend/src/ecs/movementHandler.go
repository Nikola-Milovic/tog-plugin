package ecs

import "github.com/Nikola-Milovic/tog-plugin/src/action"

type MovementHandler struct {
	manager *EntityManager
}

func (h *MovementHandler) HandleAction(indx int) {
	action, ok := h.manager.Actions[indx].(action.MovementAction)

	if !ok {
		println("error")
	}

	destination := action.Destination

	direction := h.manager.PositionComponents[indx].Position.Subtract(destination)

	println(direction)

}
