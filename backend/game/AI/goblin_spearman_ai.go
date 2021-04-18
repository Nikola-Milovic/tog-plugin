package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/game"
)

type GoblinSpearmanAI struct {
	World *game.World
}

func (ai GoblinSpearmanAI) PerformAI(index int) {
	// //Check if an enemy is in range or move to somewhere
	// closestIndex := -1
	// closestDistance := 100000
	// for _, indx := range nearbyEntities {
	// 	if w.EntityManager.GetEntities()[index].PlayerTag != w.EntityManager.GetEntities()[indx].PlayerTag {
	// 		tarPos := w.ObjectPool.Components["PositionComponent"][indx].(components.PositionComponent)

	// 		distToTarget := w.Grid.GetDistance(tarPos.Position, posComp.Position)

	// 		// If target in range, throw spear, if not check if we can attack it
	// 		if (tarPos.Position.X == posComp.Position.X || tarPos.Position.Y == tarPos.Position.Y) && distToTarget > atkComp.Range && distToTarget <= 8 && canActivateAbility(abComp.Abilities["ab_spear_throw"]["last_activated"].(int), "ab_spear_throw", w) {
	// 			fmt.Printf("Throwing spear at distance of %v", distToTarget)
	// 			data := make(map[string]interface{}, 3)
	// 			data["emitter"] = id
	// 			data["abilityID"] = "ab_spear_throw"
	// 			data["where"] = tarPos.Position
	// 			ev := engine.Event{Index: index, ID: constants.LineShotAbilityEvent, Priority: constants.SummonAbilityEventPriority, Data: data}
	// 			abComp.Abilities["ab_spear_throw"]["last_activated"] = w.Tick
	// 			w.ObjectPool.Components["AbilitiesComponent"][index] = abComp
	// 			w.EventManager.SendEvent(ev)
	// 			return
	// 		}
}
