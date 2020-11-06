package ecs

import (
	"github.com/Nikola-Milovic/tog-plugin/src/action"
)

//AttackHandler is a handler used to handle Attacking, WIP
type AttackHandler struct {
	manager *EntityManager
}

//HandleAction handles Attack Action for entity at the given index
func (h AttackHandler) HandleAction(act action.Action) {
	// action, ok := h.manager.Actions[index].(action.AttackAction)

	// if !ok {
	// 	fmt.Println("Error")
	// }

	//	fmt.Printf("I at attacking %v \n)
	//	fmt.Println(h.manager.PositionComponents[index].Position)
}
