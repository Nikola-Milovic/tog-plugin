package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

type KnightAI struct {
	World *game.World
}

func (ai KnightAI) PerformAI(index int) {
	w := ai.World

	atkComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	posComp := w.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	movComp := w.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)
	abComp := w.ObjectPool.Components["AbilitiesComponent"][index].(components.AbilitiesComponent)

	canAttack := false
	canMove := false

	if w.Tick-atkComp.TimeSinceLastAttack >= atkComp.AttackSpeed {
		canAttack = true
	} else {
		//fmt.Printf("Cannot attack %v, on cooldown for %v\n", index, atkComp.AttackSpeed-(w.Tick-atkComp.TimeSinceLastAttack))
	}

	if w.Tick-movComp.TimeSinceLastMovement >= movComp.MovementSpeed {
		canMove = true
	} else {
		//	fmt.Printf("Cannot move %v, on cooldown for %v\n", index, movComp.MovementSpeed-(w.Tick-movComp.TimeSinceLastMovement))
	}

	if !canAttack && !canMove {
		//	fmt.Printf("Cant move or attack %v\n", index)
		return
	}

	nearbyEntities := game.GetNearbyEntities(40, w, index)

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
				if canAttack {
					data := make(map[string]interface{}, 1)
					data["target"] = w.EntityManager.IndexMap[atkComp.Target]
					ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
					w.EventManager.SendEvent(ev)
					return
				}
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
			if w.Grid.GetDistanceIncludingDiagonal(tarPos.Position, posComp.Position) <= atkComp.Range {
				if canAttack {
					data := make(map[string]interface{}, 1)
					data["target"] = indx
					ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
					w.EventManager.SendEvent(ev)
					return
				}

				return
			}

			//Cast Spell
			if canActivateAbility(abComp.Ability.LastActivated, abComp.Ability.AbilityID, w) && w.Grid.GetDistanceIncludingDiagonal(tarPos.Position, posComp.Position) <= atkComp.Range+2 {
				data := make(map[string]interface{}, 1)
				data["target"] = indx
				data["abilityID"] = abComp.Ability.AbilityID
				ev := engine.Event{Index: index, ID: constants.AbilityCastEvent, Priority: constants.AbilityCastEventPriority, Data: data}
				abComp.Ability.LastActivated = w.Tick
				w.ObjectPool.Components["AbilitiesComponent"][index] = abComp

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

	if canAttack {
		//Reset target to noone
		atkComp.Target = ""
		w.ObjectPool.Components["AttackComponent"][index] = atkComp
	}
	if !canMove {
		return
	}

	if closestIndex == -1 {
		return
	}

	data := make(map[string]interface{}, 1)
	data["target"] = closestIndex
	ev := engine.Event{Index: index, ID: "MovementEvent", Priority: 99, Data: data}
	w.EventManager.SendEvent(ev)
}
