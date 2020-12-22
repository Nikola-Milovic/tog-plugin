package components

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type AttackComponent struct {
	Target              int
	Damage              int
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

	attackSpeed := int(compData[constants.AttackSpeedJson].(float64))
	damage := int(compData[constants.DamageJson].(float64))

	component.Damage = damage
	component.AttackSpeed = attackSpeed

	return component
}
