package handlers

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//ApplyEffectEventHandler is a handler
type ApplyEffectEventHandler struct {
	World *game.World
}

//HandleEvent handles
func (h ApplyEffectEventHandler) HandleEvent(ev engine.Event) {

	effID := ev.Data["effectID"].(string)

	eff := h.World.EffectDataMap[effID].(map[string]interface{})

	switch eff["Type"].(string) {
	case "dot_effect":
		{
			target := ev.Data["target"].(int)
			effect := make(map[string]interface{})
			effect["effectID"] = effID
			effect["type"] = "dot_effect"
			effect["expires"] = h.World.Tick + int(eff["Duration"].(float64))
			effect["lastTicked"] = h.World.Tick
			effComp := h.World.ObjectPool.Components["EffectsComponent"][target].(components.EffectsComponent)
			effComp.Effects = append(effComp.Effects, effect)
			h.World.ObjectPool.Components["EffectsComponent"][target] = effComp
		}
	}

}

//event := engine.Event{}
// event.ID = "ApplyDOTEffectEvent"
// event.Index = ev.Index
// event.Priority = 100
// event.Data = ev.Data
// h.World.EventManager.SendEvent(event)
