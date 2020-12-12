package game

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
)

//MovementAction specifies that the entity will be moving this tick to a given position
type MovementAction struct {
	State    string
	priority int
	Index    int
	Target   int
}

func (a MovementAction) GetPriority() int {
	return 3
}

func (a MovementAction) GetActionType() string {
	return constants.ActionTypeMovement
}

//AttackAction specifies that the entity will be attacking this round
type AttackAction struct {
	State    string
	priority int
	Index    int
	Target   int
}

func (a AttackAction) GetPriority() int {
	return 5
}

func (a AttackAction) GetActionType() string {
	return "attack"
}

//EmptyAction is a spinoff of Null Object pattern, we will use this empty action for entities who aren't doing anything this tick
type EmptyAction struct {
	Index    int
	priority int
}

func (a EmptyAction) GetPriority() int {
	return -10000
}

func (a EmptyAction) GetActionType() string {
	return constants.ActionTypeEmpty
}
