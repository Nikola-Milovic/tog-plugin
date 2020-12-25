package components

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type AbilitiesComponent struct {
	Ability engine.Ability
}

func (a AbilitiesComponent) ComponentName() string {
	return "AbilitiesComponent"
}

func AbilitiesComponentMaker(data interface{}) engine.Component {
	component := AbilitiesComponent{}

	ability := engine.Ability{}

	ability.AbilityID = data.(map[string]interface{})["AbilityID"].(string)

	component.Ability = ability

	return component
}
