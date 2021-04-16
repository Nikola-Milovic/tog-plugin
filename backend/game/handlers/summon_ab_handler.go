package handlers

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

//SummonAbilityEventHandler is a handler
type SummonAbilityEventHandler struct {
	World *game.World
}

//HandleEvent handles
//TODO: add summon bounding box
func (h SummonAbilityEventHandler) HandleEvent(ev engine.Event) {
	// abilityID := ev.Data["abilityID"].(string)
	// abilityData := h.World.AbilityDataMap[abilityID]

	// count := int(abilityData["Count"].(float64))
	// posComp := h.World.ObjectPool.Components["PositionComponent"][ev.Index].(components.PositionComponent)

	// for i := 0; i < count; i++ {
	// 	unitData := h.World.UnitDataMap[abilityData["Summon"].(string)]
	// 	where := GetClosestFreeTile(h.World, posComp.Position, posComp.Position, posComp.BoundingBox, engine.Vector{10, 10})

	// 	caster := h.World.EntityManager.GetEntities()[ev.Index]

	// 	unit := engine.NewEntityData{
	// 		PlayerTag: caster.PlayerTag,
	// 		Data:      unitData,
	// 		ID:        abilityData["Summon"].(string),
	// 		Position:  where,
	// 	}

	// 	summonIndex, summonID := h.World.EntityManager.AddEntity(unit, caster.PlayerTag, false)

	// 	//h.World.Grid.OccupyCells(where, summonID, engine.Vector{10, 10}) //TODO

	// 	pos := h.World.ObjectPool.Components["PositionComponent"][summonIndex].(components.PositionComponent)
	// 	pos.Position = where

	// 	h.World.Players[caster.PlayerTag].NumberOfUnits++
	// 	//Client event

	// 	data := make(map[string]interface{}, 5)

	// 	data["event"] = "summon"
	// 	data["tag"] = caster.PlayerTag
	// 	data["where"] = where
	// 	data["what"] = abilityData["Summon"]
	// 	data["id"] = summonID

	// 	h.World.ClientEventManager.AddEvent(data)
	// }
}
