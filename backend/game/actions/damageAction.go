package actions

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

//data["target"] = collisionID
//data["data"] = action
func DamageAction(data map[string]interface{}, w *game.World) {

	actionData := data["action_data"].(map[string]interface{}) //data from the json of the component

	ev := engine.Event{}
	evData := make(map[string]interface{}, 5)
	evData["target"] = data["target"]
	evData["amount"] = int(actionData["Damage"].(float64))
	evData["type"] = "physical"
	ev.Data = evData
	ev.ID = constants.TakeDamageEvent
	ev.Priority = constants.TakeDamageEventPriority
	ev.Index = -1 // todo
	actionData = nil

	w.EventManager.SendEvent(ev)
}
