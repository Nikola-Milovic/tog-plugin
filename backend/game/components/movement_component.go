package components

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/math"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type MovementComponent struct {
	MovementSpeed float32
	Velocity      math.Vector
	Acceleration  math.Vector
	Goal          math.Vector
	GoalMultiplier float32
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
	speed := compData[constants.MovementSpeedJson].(string)
	switch speed {
	case "slow":
		component.MovementSpeed = constants.MovementSpeedSlow
	case "medium":
		component.MovementSpeed = constants.MovementSpeedMedium
	case "fast":
		component.MovementSpeed = constants.MovementSpeedFast
	}

	//component.Velocity = math.Vector{0.1, 0.1}

	component.GoalMultiplier = 1.0

	return component
}
