package actions

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

func DamageAction(data map[string]interface{}, w *game.World) {
	ev := engine.Event{}
	data["index"] = data["target"]
	data["amount"] = data["damage"]
	data["type"] = "physical"
	ev.Data = data
	ev.ID = constants.TakeDamageEvent
	ev.Priority = constants.TakeDamageEventPriority
	ev.Index = -1 // todo

	w.EventManager.SendEvent(ev)
}
