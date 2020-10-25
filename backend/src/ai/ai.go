package ai

import (
	"backend/src/action"
	"backend/src/constants"
)

type AI interface {
	calculateAction(uindex int8) action.Action
}

type KnightAI struct{}

func (ai *KnightAI) calculateAction(uindex int8) action.Action {
	return &action.MovementAction{Destination: constants.V2{X: 50, Y: 50}}
}
