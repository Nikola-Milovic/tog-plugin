package game

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type KnightAI struct{}

func (ai KnightAI) CalculateAction(index int, e *engine.EntityManager) engine.Action {
	return EmptyAction{}
}
