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
	if atkComp.IsAttacking {
		return
	}

	adjustedAtkRange := atkComp.Range*16 + posComp.BoundingBox.DivideScalar(2).X + 2

	entities := w.EntityManager.GetEntities()

	opposingTag := 0
	if entities[index].PlayerTag == 0 {
		opposingTag = 1
	}
	nearbyEntities := helper.GetTaggedNearbyEntities(adjustedAtkRange, w, index, opposingTag)

	if len(nearbyEntities) > 0 {
		data := make(map[string]interface{}, 2)
		data["target"] = entities[nearbyEntities[0]].ID
		data["emitter"] = entities[index].ID
		ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
		ai.sendNonMovementEvent(ev, index)
		return
	}

	//Reset target to noone
	atkComp.Target = ""
	w.ObjectPool.Components["AttackComponent"][index] = atkComp

	data := make(map[string]interface{}, 2)
	data["emitter"] = entities[index].ID
	data["tag"] = entities[index].PlayerTag
	ev := engine.Event{Index: index, ID: "MovementEvent", Priority: 99, Data: data}
	w.EventManager.SendEvent(ev)
}

func (ai GenericAI) sendNonMovementEvent(ev engine.Event, index int) {
	movComp := ai.World.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)
	movComp.IsMoving = false
	movComp.Velocity = engine.Zero()
	movComp.Velocity = engine.Zero()
	ai.World.ObjectPool.Components["MovementComponent"][index] = movComp

	ai.World.EventManager.SendEvent(ev)
}
