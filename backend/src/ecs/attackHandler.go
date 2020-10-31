package ecs

//AttackHandler is a handler used to handle Attacking, WIP
type AttackHandler struct {
	manager *EntityManager
}

//HandleAction handles Attack Action for entity at the given index
func (h AttackHandler) HandleAction(index int) {
	// action, ok := h.manager.Actions[index].(action.AttackAction)

	// if !ok {
	// 	fmt.Println("Error")
	// }

	//fmt.Printf("I at %v is attacking %v \n", index, action.Target)
	//	fmt.Println(h.manager.PositionComponents[index].Position)
}
