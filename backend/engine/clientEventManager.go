package engine

//ClientEventManager holds the relevant events client needs to reconstruct the state of the game,
//alongside it also holds converters that convert the engine event to a client friendly event. Ie turns
//movementEvent into "ENTITY XYZ WALKS TO X SPOT"
type ClientEventManager struct {
	Converters map[string]ClientEventConverter
	Events     []map[string]interface{}
}

//ClientEventConverter convert events to a map of data used by the client to reconstruct the event
type ClientEventConverter func(Event) map[string]interface{}

//OnEvent dispatches the event to corresponding converter and adds it to the slice
func (ce *ClientEventManager) OnEvent(ev Event) {
	if converter, ok := ce.Converters[ev.ID]; ok {
		ce.Events = append(ce.Events, converter(ev))
	}
}

//RegisterClientEventConverter registers a new event converter, the key to the converter is the EventID of the event it converts
func (ce *ClientEventManager) RegisterClientEventConverter(conv ClientEventConverter, eventID string) {
	ce.Converters[eventID] = conv
}

//CreateClientEventManager creates and initailizes ClientEventManager
func CreateClientEventManager() *ClientEventManager {
	manager := ClientEventManager{}

	manager.Converters = make(map[string]ClientEventConverter, 10)
	manager.Events = make([]map[string]interface{}, 0, 50)

	return &manager
}
