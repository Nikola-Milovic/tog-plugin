package ecs

import (
	"github.com/Nikola-Milovic/tog-plugin/src/ai"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

type MovementComponent struct {
	Speed rune
}

type PositionComponent struct {
	X        float32
	Position constants.V2
}

type AttackComponent struct {
	AttackType string
}

type AIComponent struct {
	AI ai.AI
}
