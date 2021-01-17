package handlers

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//LineshotAbilityEventHandler is a handler
type LineshotAbilityEventHandler struct {
	World *game.World
}

//	data["emitter"] = id
// data["abilityID"] = "ab_spear_throw"
// data["where"] = tarPos.Position
//HandleEvent handles
func (h LineshotAbilityEventHandler) HandleEvent(ev engine.Event) {
	abilityID := ev.Data["abilityID"].(string)
	casterID := ev.Data["emitter"].(string)
	targetPos := ev.Data["where"].(engine.Vector)
	abData := h.World.AbilityDataMap[abilityID]

	data := make(map[string]interface{}, 5)

	data["caster"] = casterID
	data["speed"] = abData["Speed"].(int)
	data["range"] = abData["Range"].(int)

	target := engine.Vector{}
	casterPos := h.World.ObjectPool.Components["PositionComponent"][h.World.EntityManager.IndexMap[casterID]].(components.PositionComponent)

	if casterPos.Position.Y == targetPos.Y {
		target.Y = casterPos.Position.Y
		if casterPos.Position.X > targetPos.X {
			target.X = int(engine.Max(0, casterPos.Position.X-abData["Range"].(int)))
		} else {
			//Todo check the size 800/32 and make it constant
			target.X = int(engine.Min(800/32+1, casterPos.Position.X+abData["Range"].(int)))
		}
	}

	if casterPos.Position.X == targetPos.X {
		target.X = casterPos.Position.X
		if casterPos.Position.Y > targetPos.Y {
			target.Y = int(engine.Max(0, casterPos.Position.Y-abData["Range"].(int)))
		} else {
			//Todo check the size 800/32 and make it constant
			target.Y = int(engine.Min(500/32+1, casterPos.Position.Y+abData["Range"].(int)))
		}
	}

	data["target"] = target

	h.World.EntityManager.AddTempSystem("LineshotTempSystem", data)
}
