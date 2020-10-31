package action

import "github.com/Nikola-Milovic/tog-plugin/src/constants"

//Action represents a single activity/ action that the entity will perform that Tick, AI decides on the best action that the
//given entity should perform
type Action interface {
	GetPriority() rune
	GetActionState() string
}

//MovementAction specifies that the entity will be moving this tick to a given position
type MovementAction struct {
	State    string
	priority rune
	Target   constants.V2
}

func (a MovementAction) GetPriority() rune {
	return 3
}

func (a MovementAction) GetActionState() string {
	return "walk"
}

//AttackAction specifies that the entity will be attacking this round
type AttackAction struct {
	State    string
	priority rune
	Target   int
}

func (a AttackAction) GetPriority() rune {
	return 5
}

func (a AttackAction) GetActionState() string {
	return "attack"
}

//EmptyAction is a spinoff of Null Object pattern, we will use this empty action for entities who aren't doing anything this tick
type EmptyAction struct {
	priority rune
}

func (a EmptyAction) GetPriority() rune {
	return -10000
}

func (a EmptyAction) GetActionState() string {
	return "empty"
}
