package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

func canActivateAbility(lastActivated int, abilityID string, w *game.World) bool {

	if w.Tick-lastActivated > int(constants.AbilityDataMap[abilityID]["Cooldown"].(float64)) {
		return true
	}

	return false
}

func resetMovementIfNotMoving(ev engine.Event, index int, w *game.World) {
	movComp := w.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)
	movComp.IsMoving = false
	w.ObjectPool.Components["MovementComponent"][index] = movComp
}

func attackTarget(index int, target string, id string, w *game.World) {
	attackComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	attackComp.Target = target
	attackComp.TimeSinceLastAttack = w.Tick
	attackComp.IsAttacking = true
	w.ObjectPool.Components["AttackComponent"][index] = attackComp

	clientEvent := make(map[string]interface{}, 3)
	clientEvent["event"] = "attack"
	clientEvent["me"] = id
	clientEvent["who"] = target
	w.ClientEventManager.AddEvent(clientEvent)
}

func moveTowardsPoint(index int, x, y int, w *game.World) {
	movementComp := w.GetObjectPool().Components["MovementComponent"][index].(components.MovementComponent)
	positionComp := w.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)

	movementComp.IsMoving = true
	movementComp.DesiredDirection = engine.Vector{X: float32(x * constants.TileSize), Y: float32(y * constants.TileSize)}.Subtract(positionComp.Position).Normalize()

	w.GetObjectPool().Components["MovementComponent"][index] = movementComp
}
