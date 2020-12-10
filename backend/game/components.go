package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type MovementComponent struct {
	MovementSpeed         int
	TimeSinceLastMovement int
	Path                  []engine.Vector
	Target                engine.Vector
}

func (m MovementComponent) ComponentName() string {
	return "MovementComponent"
}

func MovementComponentMaker(data interface{}) engine.Component {
	compData, ok := data.(map[string]interface{})
	if !ok {
		panic(fmt.Sprint("Data given to component isn't correct type, expected map[string]interface{}"))
	}

	component := MovementComponent{}
	component.Target = engine.Vector{X: -1, Y: -1}
	speed, ok := compData[constants.MovementSpeedJson].(string)
	if !ok {
		speed = "slow"
	}

	switch speed {
	case "slow":
		component.MovementSpeed = constants.MovementSpeedSlow
	case "fast":
		component.MovementSpeed = constants.MovementSpeedFast
	}

	return component
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
	Target              int
	AttackSpeed         int
	TimeSinceLastAttack int
	Range               int
}

func (a AttackComponent) ComponentName() string {
	return "AttackComponent"
}

func AttackComponentMaker(data interface{}) engine.Component {

	compData, ok := data.(map[string]interface{})

	if !ok {
		panic(fmt.Sprint("Data given to attack component isn't correct type, expected map[string]interface{}"))
	}

	component := AttackComponent{Target: -1}

	speed := compData[constants.AttackSpeedJson]

	switch speed {
	case "slow":
		component.AttackSpeed = constants.AttackSpeedFast
	case "fast":
		component.AttackSpeed = constants.AttackSpeedSlow
	}

	return component
}

//AIComponent is used to store the AI for the specific entity
//TODO: make this to be a pointer to the same AI, maybe ditch the AI component and just make it a slice of pointers to AI
// as same units can just share the AI no need to create mulitple
type AIComponent struct {
	AI engine.AI
}
