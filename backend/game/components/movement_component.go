package components

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type MovementComponent struct {
	MovementSpeed int
	Path          []engine.Vector
	Target        engine.Vector
	IsMoving      bool
}

func (m MovementComponent) ComponentName() string {
	return "MovementComponent"
}

func MovementComponentMaker(data interface{}, additionalData map[string]interface{}, world engine.WorldI) engine.Component {
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
	case "medium":
		component.MovementSpeed = constants.MovementSpeedMedium
	case "fast":
		component.MovementSpeed = constants.MovementSpeedFast
	}

	return component
}
