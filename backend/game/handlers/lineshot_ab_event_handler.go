package handlers

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
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
	// abilityID := ev.Data["abilityID"].(string)
	// casterID := ev.Data["emitter"].(string)
	// targetPos := ev.Data["where"].(engine.Vector)
	// abData := h.World.AbilityDataMap[abilityID]

	// data := make(map[string]interface{}, 7)

	// data["caster"] = casterID
	// data["ability_id"] = abilityID
	// data["speed"] = int(abData["Speed"].(float64))
	// abRange := int(abData["Range"].(float64))

	// clientEvent := make(map[string]interface{}, 7)

	// target := engine.Vector{}
	// casterPos := h.World.ObjectPool.Components["PositionComponent"][h.World.EntityManager.GetIndexMap()[casterID]].(components.PositionComponent)

	// // N W S E
	// if casterPos.Position.Y == targetPos.Y {
	// 	target.Y = casterPos.Position.Y
	// 	if casterPos.Position.X > targetPos.X { //  Left
	// 		target.X = engine.Max(0, casterPos.Position.X-abRange)
	// 		data["position"] = engine.Vector{Y: target.Y, X: casterPos.Position.X - 1}

	// 		clientEvent["direction"] = "left"
	// 	} else { // Right
	// 		//Todo check the size 800/32 and make it constant
	// 		target.X = engine.Min(800/32+1, casterPos.Position.X+abRange)
	// 		data["position"] = engine.Vector{Y: target.Y, X: casterPos.Position.X + 1}

	// 		clientEvent["direction"] = "right"
	// 	}
	// }

	// if casterPos.Position.X == targetPos.X {
	// 	target.X = casterPos.Position.X
	// 	if casterPos.Position.Y > targetPos.Y { // Up
	// 		target.Y = engine.Max(0, casterPos.Position.Y-abRange)
	// 		data["position"] = engine.Vector{X: target.X, Y: casterPos.Position.Y - 1}

	// 		clientEvent["direction"] = "up"
	// 	} else { //Down
	// 		//Todo check the size 800/32 and make it constant
	// 		target.Y = engine.Min(500/32+1, casterPos.Position.Y+abRange)
	// 		data["position"] = engine.Vector{X: target.X, Y: casterPos.Position.Y + 1}

	// 		clientEvent["direction"] = "down"
	// 	}
	// }

	// //DIAGONAL

	// data["target"] = target
	// data["last_moved"] = h.World.Tick
	// data["id"] = engine.MustGenerateID()
	// data["projectile"] = abData["Projectile"]

	// clientEvent["position"] = data["position"]
	// clientEvent["target"] = target
	// clientEvent["event"] = "lineshot_projectile_spawn"
	// clientEvent["id"] = data["id"]
	// clientEvent["projectile"] = abData["Projectile"]

	// h.World.EntityManager.AddTempSystem("LineshotTempSystem", data, h.World)
	// h.World.ClientEventManager.AddEvent(clientEvent)
}
