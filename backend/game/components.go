package game

import "github.com/Nikola-Milovic/tog-plugin/engine"

type MovementComponent struct {
	Tick int
	Path []engine.Vector
}

func (m MovementComponent) ComponentName() string {
	return "MovementComponent"
}

func MovementComponentMaker(data interface{}) engine.Component {
	return MovementComponent{}
}

type PositionComponent struct {
	Position engine.Vector
}

func (m PositionComponent) ComponentName() string {
	return "PositionComponent"
}

func PositionComponentMaker(data interface{}) engine.Component {
	return PositionComponent{}
}

type AttackComponent struct {
	Target int
	Type   string
	Range  int
}

//AIComponent is used to store the AI for the specific entity
//TODO: make this to be a pointer to the same AI, maybe ditch the AI component and just make it a slice of pointers to AI
// as same units can just share the AI no need to create mulitple
type AIComponent struct {
	AI engine.AI
}
