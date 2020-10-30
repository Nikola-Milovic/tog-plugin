package ecs

//AttackHandler is a handler used to handle Attacking, WIP
type AttackHandler struct {
	manager *EntityManager
}

//HandleAction handles Attack Action for entity at the given index
func (h AttackHandler) HandleAction(index int) {
	//action, ok := h.manager.Actions[index].(action.AttackAction)

	// if !ok {
	// 	fmt.Println("Error")
	// }

	//fmt.Printf("I %v is attacking %v \n", index, action.Target)
	// fmt.Println(direction.MultiplyScalar(float64(h.manager.MovementComponents[index].Speed)))
	// fmt.Println(h.manager.PositionComponents[index].Position)
}
