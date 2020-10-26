package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/src/action"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

type AI interface {
	CalculateAction(index int) action.Action
}

type KnightAI struct{}

func (ai KnightAI) CalculateAction(index int) action.Action {
	return action.MovementAction{Destination: constants.V2{X: 50, Y: 50}}
}
