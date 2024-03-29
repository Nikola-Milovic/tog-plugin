package ai

// import (
// 	"github.com/Nikola-Milovic/tog-plugin/engine"
// 	"github.com/Nikola-Milovic/tog-plugin/game"
// 	"github.com/Nikola-Milovic/tog-plugin/game/components"
// 	"github.com/Nikola-Milovic/tog-plugin/game/helper"
// )

// type KnightAI struct {
// 	World *game.World
// }

// func (ai KnightAI) PerformAI(index int) {
// 	w := ai.World

// 	atkComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
// 	posComp := w.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
// 	movComp := w.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)
// 	//abComp := w.ObjectPool.Components["AbilitiesComponent"][index].(components.AbilitiesComponent)

// 	//If we're moving or attacking just return
// 	if atkComp.IsAttacking || movComp.IsMoving {
// 		return
// 	}

// 	nearbyEntities := helper.GetNearbyEntities(40, w, index)

// 	target, ok := w.EntityManager.IndexMap[atkComp.Goal]

// 	if ok { // if our target still exists

// 		//Check if target is inactive now
// 		if !w.EntityManager.Entities[target].Active {
// 			atkComp.Goal = ""
// 		}

// 		//If we're already attacking, keep attacking
// 		if atkComp.Goal != "" {
// 			tarPos := w.ObjectPool.Components["PositionComponent"][target].(components.PositionComponent)
// 			if w.Grid.GetDistanceIncludingDiagonal(posComp.Position, tarPos.Position) <= atkComp.Range {

// 				data := make(map[string]interface{}, 2)
// 				data["emitter"] = ai.World.EntityManager.Entities[index].ID
// 				data["target"] = atkComp.Goal
// 				ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
// 				w.EventManager.SendEvent(ev)
// 				return

// 			}
// 		}
// 	}

// 	//Check if an enemy is in range or move to somewhere
// 	closestIndex := -1
// 	closestDistance := 100000
// 	for _, indx := range nearbyEntities {
// 		if w.EntityManager.Entities[index].PlayerTag != w.EntityManager.Entities[indx].PlayerTag {
// 			tarPos := w.ObjectPool.Components["PositionComponent"][indx].(components.PositionComponent)
// 			if w.Grid.GetDistanceIncludingDiagonal(tarPos.Position, posComp.Position) <= atkComp.Range {

// 				data := make(map[string]interface{}, 2)
// 				data["target"] = ai.World.EntityManager.Entities[indx].ID
// 				data["emitter"] = ai.World.EntityManager.Entities[index].ID
// 				ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
// 				w.EventManager.SendEvent(ev)
// 				return

// 			}

// 			// //Cast Spell
// 			// if canActivateAbility(abComp.Ability.LastActivated, abComp.Ability.AbilityID, w) && w.Grid.GetDistanceIncludingDiagonal(tarPos.Position, posComp.Position) <= atkComp.Range+2 {
// 			// 	data := make(map[string]interface{}, 2)
// 			// 	data["target"] = ai.World.EntityManager.Entities[indx].ID
// 			// 	data["emitter"] = ai.World.EntityManager.Entities[index].ID
// 			// 	data["abilityID"] = abComp.Ability.AbilityID
// 			// 	ev := engine.Event{Index: index, ID: constants.AbilityCastEvent, Priority: constants.AbilityCastEventPriority, Data: data}
// 			// 	abComp.Ability.LastActivated = w.Tick
// 			// 	w.ObjectPool.Components["AbilitiesComponent"][index] = abComp

// 			// 	w.EventManager.SendEvent(ev)
// 			// 	return
// 			// }

// 			dist := w.Grid.GetDistance(tarPos.Position, posComp.Position)
// 			if dist < closestDistance {
// 				closestIndex = indx
// 				closestDistance = dist
// 			}

// 		}
// 	}

// 	//Reset target to noone
// 	atkComp.Goal = ""
// 	w.ObjectPool.Components["AttackComponent"][index] = atkComp

// 	if closestIndex == -1 {
// 		return
// 	}

// 	data := make(map[string]interface{}, 2)
// 	data["target"] = closestIndex
// 	data["emitter"] = ai.World.EntityManager.Entities[index].ID
// 	ev := engine.Event{Index: index, ID: "MovementEvent", Priority: 99, Data: data}
// 	w.EventManager.SendEvent(ev)
// }
