package engine

//ClientEventManager holds the relevant events client needs to reconstruct the state of the game,
//movementEvent into "ENTITY XYZ WALKS TO X SPOT"
type ClientEventManager struct {
	Events     []map[string]interface{}
}

//AddEvent directly adds event to the even slice, this is primarly used for cases where converters are not needed and we just directly add the ClientEvent (Movement event for example)
func (ce *ClientEventManager) AddEvent(ev map[string]interface{}) {
	ce.Events = append(ce.Events, ev)
}


//CreateClientEventManager creates and initailizes ClientEventManager
func CreateClientEventManager() *ClientEventManager {
	manager := ClientEventManager{}

	manager.Events = make([]map[string]interface{}, 0, 50)

	return &manager
}
