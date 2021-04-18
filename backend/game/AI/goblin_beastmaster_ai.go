package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/game"	
)

type GoblinBeastMasterAI struct {
	World *game.World
}

func (ai GoblinBeastMasterAI) PerformAI(index int) {
	// // //Cast Spell
	// if canActivateAbility(abComp.Abilities["ab_summon_wolf"]["last_activated"].(int), "ab_summon_wolf", w) {
	// 	data := make(map[string]interface{}, 2)
	// 	data["emitter"] = id
	// 	data["abilityID"] = "ab_summon_wolf"
	// 	ev := engine.Event{Index: index, ID: constants.SummonAbilityEvent, Priority: constants.SummonAbilityEventPriority, Data: data}
	// 	abComp.Abilities["ab_summon_wolf"]["last_activated"] = w.Tick
	// 	w.ObjectPool.Components["AbilitiesComponent"][index] = abComp
	// 	w.EventManager.SendEvent(ev)
	// 	return
	// }
}
