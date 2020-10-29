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

	destination := action.Destination

	direction := destination.Subtract(h.manager.PositionComponents[index].Position).Normalize()

	h.manager.PositionComponents[index].Position = h.manager.PositionComponents[index].Position.Add((direction.MultiplyScalar(float64(h.manager.MovementComponents[index].Speed))))

	// fmt.Println(direction)
	// fmt.Println(direction.MultiplyScalar(float64(h.manager.MovementComponents[index].Speed)))
	// fmt.Println(h.manager.PositionComponents[index].Position)
}
