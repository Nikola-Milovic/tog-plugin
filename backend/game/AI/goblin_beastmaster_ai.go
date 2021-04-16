package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/game"	
)

type GoblinBeastMasterAI struct {
	World *game.World
}

func (ai GoblinBeastMasterAI) PerformAI(index int) {
	// w := ai.World

	// id := w.EntityManager.GetEntities()[index].ID

	// atkComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	// posComp := w.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	// movComp := w.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)
	// abComp := w.ObjectPool.Components["AbilitiesComponent"][index].(components.AbilitiesComponent)

	// //If we're moving or attacking just return
	// if atkComp.IsAttacking || movComp.IsMoving {
	// 	return
	// }

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

	// nearbyEntities := helper.GetNearbyEntities(40, w, index)

	// target, ok := w.EntityManager.GetIndexMap()[atkComp.Target]

	// if ok { // if our target still exists

	// 	//Check if target is inactive now
	// 	if !w.EntityManager.GetEntities()[target].Active {
	// 		atkComp.Target = ""
	// 	}

	// 	//If we're already attacking, keep attacking
	// 	if atkComp.Target != "" {
	// 		tarPos := w.ObjectPool.Components["PositionComponent"][target].(components.PositionComponent)
	// 		if w.Grid.GetDistanceIncludingDiagonal(posComp.Position, tarPos.Position) <= atkComp.Range {

	// 			data := make(map[string]interface{}, 2)
	// 			data["emitter"] = id
	// 			data["target"] = atkComp.Target
	// 			ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
	// 			w.EventManager.SendEvent(ev)
	// 			return

	// 		}
	// 	}
	// }
	// //Check if an enemy is in range or move to somewhere
	// closestIndex := -1
	// closestDistance := 100000
	// for _, indx := range nearbyEntities {
	// 	if w.EntityManager.GetEntities()[index].PlayerTag != w.EntityManager.GetEntities()[indx].PlayerTag {
	// 		tarPos := w.ObjectPool.Components["PositionComponent"][indx].(components.PositionComponent)
	// 		if w.Grid.GetDistanceIncludingDiagonal(tarPos.Position, posComp.Position) <= atkComp.Range {

	// 			data := make(map[string]interface{}, 2)
	// 			data["target"] = ai.World.EntityManager.GetEntities()[indx].ID
	// 			data["emitter"] = id
	// 			ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
	// 			w.EventManager.SendEvent(ev)
	// 			return

	// 		}

	// 		dist := w.Grid.GetDistance(tarPos.Position, posComp.Position)
	// 		if dist < closestDistance {
	// 			closestIndex = indx
	// 			closestDistance = dist
	// 		}

	// 	}
	// }

	// //Reset target to noone
	// atkComp.Target = ""
	// w.ObjectPool.Components["AttackComponent"][index] = atkComp

	// if closestIndex == -1 {
	// 	return
	// }

	// data := make(map[string]interface{}, 2)
	// data["target"] = closestIndex
	// data["emitter"] = id
	// ev := engine.Event{Index: index, ID: "MovementEvent", Priority: 99, Data: data}
	// w.EventManager.SendEvent(ev)
}
