package client

import "github.com/Nikola-Milovic/tog-plugin/engine"

//MovementEventConverter doesn't handle MovementEvent, but rather a hacky Movement event emitted specifically for the client
//this allows for decoupling between the AI and handlers
func MovementEventConverter(ev engine.Event) map[string]interface{} {
	data := make(map[string]interface{}, 3)

	data["event"] = "walk"
	data["who"] = ev.Data["emitter"]
	data["where"] = ev.Data["where"]

	return data
}

//	data := make(map[string]interface{}, 1)
// data["target"] = closestIndex
// data["unit"] = ai.World.EntityManager.Entities[index].ID
// ev := engine.Event{Index: index, ID: "MovementEvent", Priority: 99, Data: data}
