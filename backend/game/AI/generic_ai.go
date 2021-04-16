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

	////If we're moving or attacking just return
	//if atkComp.IsAttacking || movComp.IsMoving {
	//	return
	//}

	adjustedAtkRange := atkComp.Range*16 + posComp.BoundingBox.DivideScalar(2).X
	target, ok := w.EntityManager.GetIndexMap()[atkComp.Target]

	entities := w.EntityManager.GetEntities()

	if ok { // if our target still exists
		//Check if target is inactive now
		if !entities[target].Active {
			atkComp.Target = ""
		}

		//If we're already attacking, keep attacking
		if atkComp.Target != "" {
			tarPos := w.ObjectPool.Components["PositionComponent"][target].(components.PositionComponent)
			if engine.GetDistanceIncludingDiagonal(posComp.Position, tarPos.Position) <= adjustedAtkRange {
				data := make(map[string]interface{}, 2)
				data["emitter"] = entities[index].ID
				data["target"] = atkComp.Target
				ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
				ai.sendNonMovementEvent(ev, index)
				return
			}
		}
	}

	nearbyEntities := helper.GetNearbyEntities(adjustedAtkRange, w, index)

	for _, indx := range nearbyEntities {
		if entities[index].PlayerTag != entities[indx].PlayerTag {
			tarPos := w.ObjectPool.Components["PositionComponent"][indx].(components.PositionComponent)
			if engine.GetDistanceIncludingDiagonal(tarPos.Position, posComp.Position) <= adjustedAtkRange {

				data := make(map[string]interface{}, 2)
				data["target"] = entities[indx].ID
				data["emitter"] = entities[index].ID
				ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
				ai.sendNonMovementEvent(ev, index)
				return

			}
		}
	}

	//Reset target to noone
	atkComp.Target = ""
	w.ObjectPool.Components["AttackComponent"][index] = atkComp

	data := make(map[string]interface{}, 2)
	data["emitter"] = entities[index].ID
	ev := engine.Event{Index: index, ID: "MovementEvent", Priority: 99, Data: data}
	w.EventManager.SendEvent(ev)
}

func (ai GenericAI) sendNonMovementEvent(ev engine.Event, index int) {
	movComp := ai.World.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)
	movComp.IsMoving = false
	movComp.Direction = engine.Zero()
	movComp.Direction = engine.Zero()
	ai.World.ObjectPool.Components["MovementComponent"][index] = movComp

	ai.World.EventManager.SendEvent(ev)
}