package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/helper"
)

func canActivateAbility(lastActivated int, abilityID string, w *game.World) bool {

	if w.Tick-lastActivated > int(constants.AbilityDataMap[abilityID]["Cooldown"].(float64)) {
		return true
	}

	return false
}

func attackTarget(index int, target int, id int, w *game.World) {
	attackComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	attackComp.Target = target
	attackComp.TimeSinceLastAttack = w.Tick
	w.ObjectPool.Components["AttackComponent"][index] = attackComp

	helper.SwitchState(w.EntityManager.GetEntities(), index, constants.StateAttacking, w)

	clientEvent := make(map[string]interface{}, 3)
	clientEvent["event"] = "attack"
	clientEvent["me"] = id
	clientEvent["who"] = target
	w.ClientEventManager.AddEvent(clientEvent)
}

func moveTowardsTarget(index int, targetID int, w *game.World, isEngaging bool) {
	movementComp := w.GetObjectPool().Components["MovementComponent"][index].(components.MovementComponent)

	if isEngaging {
		helper.SwitchState(w.EntityManager.GetEntities(), index, constants.StateEngaging, w)
	} else {
		helper.SwitchState(w.EntityManager.GetEntities(), index, constants.StateWalking, w)
	}
	movementComp.TargetID = targetID

	w.GetObjectPool().Components["MovementComponent"][index] = movementComp
}

func getEnemyTag(tag int) int {
	if tag == 0 {
		return 1
	} else {
		return 0
	}
}