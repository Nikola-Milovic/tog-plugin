package action

import "github.com/Nikola-Milovic/tog-plugin/src/constants"

type Action interface {
	GetPriority() rune
}

type MovementAction struct {
	priority    rune
	Destination constants.V2
}

func (a MovementAction) GetPriority() rune {
	return 10
}

type EmptyAction struct {
	priority rune
}

func (a EmptyAction) GetPriority() rune {
	return -10000
}
