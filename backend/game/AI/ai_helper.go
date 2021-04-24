package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

func canActivateAbility(lastActivated int, abilityID string, w *game.World) bool {

	if w.Tick-lastActivated > int(constants.AbilityDataMap[abilityID]["Cooldown"].(float64)) {
		return true
	}

	return false
}

func attackTarget(index int, target string, id string, w *game.World) {
	attackComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	attackComp.Target = target
	attackComp.TimeSinceLastAttack = w.Tick
	w.ObjectPool.Components["AttackComponent"][index] = attackComp

	movementComp := w.GetObjectPool().Components["MovementComponent"][index].(components.MovementComponent)

	w.EntityManager.GetEntities()[index].State = constants.StateAttacking

	w.GetObjectPool().Components["MovementComponent"][index] = movementComp

	clientEvent := make(map[string]interface{}, 3)
	clientEvent["event"] = "attack"
	clientEvent["me"] = id
	clientEvent["who"] = target
	w.ClientEventManager.AddEvent(clientEvent)
}

func moveTowardsTarget(index int, targetID string, w *game.World, isEngaging bool) {
	movementComp := w.GetObjectPool().Components["MovementComponent"][index].(components.MovementComponent)

	if isEngaging {
		w.EntityManager.GetEntities()[index].State = constants.StateEngaging
	} else {
		w.EntityManager.GetEntities()[index].State = constants.StateWalking
	}
	movementComp.TargetID = targetID

	w.GetObjectPool().Components["MovementComponent"][index] = movementComp
}
