package handlers

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

//AbilityCastEventHandler is a handler
type AbilityCastEventHandler struct {
	World *game.World
}

//HandleEvent handles
func (h AbilityCastEventHandler) HandleEvent(ev engine.Event) {
	abilityID := ev.Data["abilityID"].(string)
	switch h.World.AbilityDataMap[abilityID]["Type"] {
	case "singleTarget":
		{
			singleTarEvent := engine.Event{}
			singleTarEvent.ID = constants.SingleTargetAbilityEvent
			singleTarEvent.Index = ev.Index
			singleTarEvent.Priority = constants.SingleTargetAbilityEventPriority
			singleTarEvent.Data = ev.Data

			h.World.EventManager.SendEvent(singleTarEvent)

		}
	}

}
