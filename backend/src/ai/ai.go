package ai

import (
	"backend/src/action"
	"backend/src/constants"
)

type AI interface {
	CalculateAction(uindex uint16) action.Action
}

type KnightAI struct{}

func (ai *KnightAI) CalculateAction(uindex uint16) action.Action {
	return &action.MovementAction{Destination: constants.V2{X: 50, Y: 50}}
}
