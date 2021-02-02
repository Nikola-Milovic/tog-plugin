package ai

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
)

type GoblinSpearmanAI struct {
	World *game.World
}

func (ai GoblinSpearmanAI) PerformAI(index int) {
	w := ai.World

	id := w.EntityManager.Entities[index].ID

	atkComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	posComp := w.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	movComp := w.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)
	abComp := w.ObjectPool.UniqueComponents["AbilitiesComponent"][id].(components.AbilitiesComponent)

	//If we're moving or attacking just return
	if atkComp.IsAttacking || movComp.IsMoving {
		return
	}

	nearbyEntities := helper.GetNearbyEntities(40, w, index)

	target, ok := w.EntityManager.IndexMap[atkComp.Target]

	if ok { // if our target still exists

		//Check if target is inactive now
		if !w.EntityManager.Entities[target].Active {
			atkComp.Target = ""
		}

		//If we're already attacking, keep attacking
		if atkComp.Target != "" {
			tarPos := w.ObjectPool.Components["PositionComponent"][target].(components.PositionComponent)
			if w.Grid.GetDistanceIncludingDiagonal(posComp.Position, tarPos.Position) <= atkComp.Range {

				data := make(map[string]interface{}, 2)
				data["emitter"] = id
				data["target"] = atkComp.Target
				ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
				w.EventManager.SendEvent(ev)
				return

			}
		}
	}
	//Check if an enemy is in range or move to somewhere
	closestIndex := -1
	closestDistance := 100000
	for _, indx := range nearbyEntities {
		if w.EntityManager.Entities[index].PlayerTag != w.EntityManager.Entities[indx].PlayerTag {
			tarPos := w.ObjectPool.Components["PositionComponent"][indx].(components.PositionComponent)

			distToTarget := w.Grid.GetDistance(tarPos.Position, posComp.Position)

			// If target in range, throw spear, if not check if we can attack it
			if (tarPos.Position.X == posComp.Position.X || tarPos.Position.Y == tarPos.Position.Y) && distToTarget > atkComp.Range && distToTarget <= 8 && canActivateAbility(abComp.Abilities["ab_spear_throw"]["last_activated"].(int), "ab_spear_throw", w) {
				fmt.Printf("Throwing spear at distance of %v", distToTarget)
				data := make(map[string]interface{}, 3)
				data["emitter"] = id
				data["abilityID"] = "ab_spear_throw"
				data["where"] = tarPos.Position
				ev := engine.Event{Index: index, ID: constants.LineShotAbilityEvent, Priority: constants.SummonAbilityEventPriority, Data: data}
				abComp.Abilities["ab_spear_throw"]["last_activated"] = w.Tick
				w.ObjectPool.UniqueComponents["AbilitiesComponent"][id] = abComp
				w.EventManager.SendEvent(ev)
				return
			}

			if distToTarget <= atkComp.Range {
				data := make(map[string]interface{}, 2)
				data["target"] = ai.World.EntityManager.Entities[indx].ID
				data["emitter"] = id
				ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
				w.EventManager.SendEvent(ev)
				return

			}

			dist := w.Grid.GetDistance(tarPos.Position, posComp.Position)
			if dist < closestDistance {
				closestIndex = indx
				closestDistance = dist
			}

		}
	}

	//Reset target to noone
	atkComp.Target = ""
	w.ObjectPool.Components["AttackComponent"][index] = atkComp

	if closestIndex == -1 {
		return
	}

	data := make(map[string]interface{}, 2)
	data["target"] = closestIndex
	data["emitter"] = id
	ev := engine.Event{Index: index, ID: "MovementEvent", Priority: 99, Data: data}
	w.EventManager.SendEvent(ev)
}
