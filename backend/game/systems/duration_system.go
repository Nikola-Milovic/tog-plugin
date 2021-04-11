package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

type DurationSystem struct {
	World *game.World
}

func (ds DurationSystem) Update() {
	for compIndex, comp := range ds.World.ObjectPool.Components["EffectsComponent"] {
		for index, eff := range comp.(components.EffectsComponent).Effects {
			if expiresOnTick, ok := eff["expires"]; ok {
				if ds.World.Tick > expiresOnTick.(int) {
					println("Effect expired")
					component := comp.(components.EffectsComponent)
					component.Effects = engine.RemoveFromSliceMapStringInterface(comp.(components.EffectsComponent).Effects, index)
					ds.World.ObjectPool.Components["EffectsComponent"][compIndex] = component
				}
			}

		}
	}
}
