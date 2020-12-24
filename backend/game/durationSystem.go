package game

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

type DurationSystem struct {
	world *World
}

func (ds DurationSystem) Update() {
	for compIndex, comp := range ds.world.ObjectPool.Components["EffectsComponent"] {
		for index, eff := range comp.(components.EffectsComponent).Effects {
			if expiresOnTick, ok := eff["expires"]; ok {
				if ds.world.Tick > expiresOnTick.(int) {
					println("Effect expired")
					component := comp.(components.EffectsComponent)
					component.Effects = engine.RemoveFromSliceMapStringInterface(comp.(components.EffectsComponent).Effects, index)
					ds.world.ObjectPool.Components["EffectsComponent"][compIndex] = component
				}
			}

		}
	}
}
