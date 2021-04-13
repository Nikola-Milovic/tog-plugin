package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
)

type GenericAI struct {
	World *game.World
}

func (ai GenericAI) PerformAI(index int) {
	w := ai.World

	atkComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	posComp := w.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	movComp := w.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)

	//If we're moving or attacking just return
	if atkComp.IsAttacking || movComp.IsMoving {
		return
	}

	nearbyEntities := helper.GetNearbyEntities(400, w, index)

	target, ok := w.EntityManager.GetIndexMap()[atkComp.Target]

	entities := w.EntityManager.GetEntities()

	adjustedAtkRange := 1 + (atkComp.Range*4 + posComp.BoundingBox.X/8) // (bounding box / 2 for half of the box and that divided for 4 for tiles)

	if ok { // if our target still exists

		//Check if target is inactive now
		if !entities[target].Active {
			atkComp.Target = ""
		}

		//If we're already attacking, keep attacking
		if atkComp.Target != "" {
			tarPos := w.ObjectPool.Components["PositionComponent"][target].(components.PositionComponent)
			if w.Grid.GetDistanceIncludingDiagonal(posComp.Position, tarPos.Position) <= adjustedAtkRange {

				data := make(map[string]interface{}, 2)
				data["emitter"] = entities[index].ID
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
		if entities[index].PlayerTag != entities[indx].PlayerTag {
			tarPos := w.ObjectPool.Components["PositionComponent"][indx].(components.PositionComponent)
			if w.Grid.GetDistanceIncludingDiagonal(tarPos.Position, posComp.Position) <= adjustedAtkRange {

				data := make(map[string]interface{}, 2)
				data["target"] = entities[indx].ID
				data["emitter"] = entities[index].ID
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
	data["emitter"] = entities[index].ID
	ev := engine.Event{Index: index, ID: "MovementEvent", Priority: 99, Data: data}
	w.EventManager.SendEvent(ev)
}
