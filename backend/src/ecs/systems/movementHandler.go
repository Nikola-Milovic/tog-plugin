package handlers

import (
	"backend/src/action"
	"backend/src/ecs"
)

type MovementHandler struct {
	EntCompSystem *ecs.ECS
}

func (h *MovementHandler) Handle(ent ecs.Entity) {
	action, ok := h.EntCompSystem.Actions[ent.Index].(action.MovementAction)

	if !ok {
		println("error")
	}

	destination := action.Destination

}
