package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

type DotSystem struct {
	world *World
}

func (ds DotSystem) Update() {
	for index, comp := range ds.world.ObjectPool.Components["EffectsComponent"] {
		for effIndex, eff := range comp.(components.EffectsComponent).Effects {
			effID := eff["effectID"].(string)
			effData := ds.world.EffectDataMap[effID].(map[string]interface{})
			if effData["Type"] == "dot_effect" {
				ticksEvery := int(effData["TicksEvery"].(float64))
				damage := int(effData["DamagePerTick"].(float64))
				//	fmt.Printf("Entity at %v is affected by dot %v, expires on %v, ticksEvery %v, damage %v\n", index, effID, eff["expires"].(int), effData["TicksEvery"].(float64), effData["DamagePerTick"].(float64))
				lastTicked := eff["lastTicked"].(int)
				if ds.world.Tick-lastTicked >= ticksEvery {
					fmt.Printf("Damage tick for dmg %v on %v\n", damage, index)
					//Set last ticked to this tick
					ds.world.ObjectPool.Components["EffectsComponent"][index].(components.EffectsComponent).Effects[effIndex]["lastTicked"] = ds.world.Tick

					//Take damage event
					takeDamageEvent := engine.Event{}
					takeDamageEvent.ID = constants.TakeDamageEvent
					takeDamageEvent.Index = index
					takeDamageEvent.Priority = constants.TakeDamageEventPriority
					data := make(map[string]interface{}, 3)
					data["index"] = index
					data["amount"] = damage
					data["type"] = "magical"
					takeDamageEvent.Data = data

					ds.world.EventManager.SendEvent(takeDamageEvent)
				}

			}
		}
	}
}
