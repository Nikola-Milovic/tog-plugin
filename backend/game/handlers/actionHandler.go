package handlers

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/actions"
)

//ActionHandler is a handler
type ActionHandler struct {
	World *game.World
}

//HandleEvent handles
func (h ActionHandler) HandleEvent(ev engine.Event) {

	actionData := ev.Data["action_data"].(map[string]interface{})
	actionID := actionData["ActionID"].(string)

	switch actionID {
	case "act_damage":
		actions.DamageAction(actionData, h.World)
	}

}
