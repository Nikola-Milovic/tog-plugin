package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

type DotSystem struct {
	world *World
}

func (ds DotSystem) Update() {
	for index, comp := range ds.world.ObjectPool.Components["EffectsComponent"] {
		for _, eff := range comp.(components.EffectsComponent).Effects {
			effID := eff["effectID"].(string)
			effData := ds.world.EffectDataMap[effID].(map[string]interface{})
			if effData["Type"] == "dot_effect" {
				fmt.Printf("Entity at %v is affected by dot %v\n", index, effID)
			}
		}
	}
}
