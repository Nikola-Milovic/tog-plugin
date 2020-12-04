package game

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type KnightAI struct {
	index       int
	targetIndex int
	state       constants.UnitState
	Active      bool
}

func (ai KnightAI) CalculateAction() engine.Action {
	return EmptyAction{}
}
