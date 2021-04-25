package helper

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

func SwitchState(entities []engine.Entity, index int, newState string, w *game.World) {
	previousState := entities[index].State

	entities[index].State = newState

	atkComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	posComp := w.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	movComp := w.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)

	switch previousState {
	case constants.StateAttacking:
		{
			atkComp.Target = -1
		}

	case constants.StateWalking:
		{
			if newState != constants.StateEngaging {
				movComp.Velocity = math.Zero()
				movComp.TargetID = -1
			}
		}
	case constants.StateEngaging:
		{
			if newState != constants.StateWalking {
				movComp.Velocity = math.Zero()
				movComp.TargetID = -1
			}
		}

	}

	w.ObjectPool.Components["AttackComponent"][index] = atkComp
	w.ObjectPool.Components["PositionComponent"][index] = posComp
	w.ObjectPool.Components["MovementComponent"][index] = movComp
}
