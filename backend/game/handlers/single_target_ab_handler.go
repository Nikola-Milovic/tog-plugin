package handlers

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/actions"
)

//SingleTargetAbilityEventHandler is a handler
type SingleTargetAbilityEventHandler struct {
	World *game.World
}

//HandleEvent handles
func (h SingleTargetAbilityEventHandler) HandleEvent(ev engine.Event) {
	fmt.Printf("Cast siingleTargetAbility by %v, and the ability is %v at %v\n",
		ev.Index, ev.Data["abilityID"].(string), ev.Data["target"].(string))

	abilityID := ev.Data["abilityID"].(string)
	abilityData := constants.AbilityDataMap[abilityID]

	switch abilityData["Action"].(map[string]interface{})["ActionID"] {
	case "act_damage":
		data := make(map[string]interface{}, 5)
		data["target"] = ev.Data["target"].(string)
		data["damage"] = int(abilityData["Action"].(map[string]interface{})["Damage"].(float64))

		actions.DamageAction(data, h.World)
	}

}
