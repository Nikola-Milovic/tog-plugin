package ecs

import (
	"backend/src/action"
)

type MovementHandler struct {
	EntCompSystem *ECS
}

func (h *MovementHandler) Handle(indx uint16) {
	action, ok := h.EntCompSystem.Actions[indx].(action.MovementAction)

	if !ok {
		println("error")
	}

	destination := action.Destination

	direction := h.EntCompSystem.PositionComponents[indx].Position.Subtract(destination)

}
