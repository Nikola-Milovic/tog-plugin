package ecs

import (
	//	"github.com/Nikola-Milovic/tog-plugin/src/ai"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

type MovementComponent struct {
	Speed        float32
	Velocity     constants.V2
	Accelaration constants.V2
}

type PositionComponent struct {
	Position constants.V2
}

type AttackComponent struct {
	Type  string
	Range int
}

//AIComponent is used to store the AI for the specific entity
//TODO: make this to be a pointer to the same AI, maybe ditch the AI component and just make it a slice of pointers to AI
// as same units can just share the AI no need to create mulitple
type AIComponent struct {
	AI AI
}
