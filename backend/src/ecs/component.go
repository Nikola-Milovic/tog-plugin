package ecs

import (
	"backend/src/ai"
	"backend/src/constants"
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
