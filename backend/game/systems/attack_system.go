package systems

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/math"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

type AttackSystem struct {
	World *game.World
}

func (as AttackSystem) Update() {
	world := as.World
	entities := world.EntityManager.GetEntities()
	indexMap := world.EntityManager.GetIndexMap()
	posComponents := world.ObjectPool.Components["PositionComponent"]
	for index, comp := range world.ObjectPool.Components["AttackComponent"] {
		attackComp := comp.(components.AttackComponent)
		return
		if entities[index].State != constants.StateAttacking {
			continue
		}

		targetIndex := indexMap[attackComp.Target]

		if as.World.Tick-attackComp.TimeSinceLastAttack == attackComp.AttackSpeed/2 {
			//Attack finished, can attack again
			target := attackComp.Target
			takeDamageEvent := engine.Event{}
			takeDamageEvent.ID = constants.TakeDamageEvent
			takeDamageEvent.Index = index
			takeDamageEvent.Priority = constants.AttackEventPriority
			data := make(map[string]interface{}, 3)
			data["target"] = target
			data["amount"] = attackComp.Damage
			data["type"] = "physical"
			takeDamageEvent.Data = data

			world.EventManager.SendEvent(takeDamageEvent)

			if attackComp.OnHit != "" {
				world.EventManager.SendEvent(onHitEvent(index, target, attackComp.OnHit))
				fmt.Printf("Send onhit event tick %v\n", world.Tick)
			}
			continue
		} else if as.World.Tick-attackComp.TimeSinceLastAttack == attackComp.AttackSpeed {
			attackComp.TimeSinceLastAttack = world.Tick

			myPos := posComponents[index].(components.PositionComponent)
			attackRange := attackComp.Range*10

			ePos := posComponents[targetIndex].(components.PositionComponent)
			dist := math.GetDistanceIncludingDiagonalVectors(myPos.Position, ePos.Position)
			if !(dist-myPos.BoundingBox.X/2-ePos.BoundingBox.X/2 <= attackRange) {
				//attackComp.IsAttacking = false
			}
		}

		if !entities[targetIndex].Active {
			//attackComp.IsAttacking = false
		}

		world.ObjectPool.Components["AttackComponent"][index] = attackComp
	}
}

func onHitEvent(index int, target string, effect string) engine.Event {
	//attackComp := h.World.ObjectPool.Components["AttackComponent"][ev.Index].(components.AttackComponent)
	event := engine.Event{}
	event.ID = constants.ApplyEffectEvent
	event.Index = index
	event.Priority = constants.ApplyEffectEventPriority
	data := make(map[string]interface{}, 1)
	data["effectID"] = effect
	data["target"] = target
	event.Data = data

	return event
}
