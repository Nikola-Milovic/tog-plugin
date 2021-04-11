package components

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type StatsComponent struct {
	MaxHealth int
	Health    int
	Armor     int
}

func (h StatsComponent) ComponentName() string {
	return "StatsComponent"
}

func StatsComponentMaker(data interface{}, additionalData map[string]interface{}, world engine.WorldI) engine.Component {

	compData, ok := data.(map[string]interface{})

	if !ok {
		panic(fmt.Sprint("Data given to stats component isn't correct type, expected map[string]interface{}"))
	}

	component := StatsComponent{}

	health := int(compData[constants.MaxHealthJson].(float64))

	component.MaxHealth = health
	component.Health = health

	return component
}
