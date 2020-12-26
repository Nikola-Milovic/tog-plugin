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

	compData := data.(map[string]interface{})

	if val, ok := compData["AbilityID"]; ok {
		ability := engine.Ability{}
		ability.AbilityID = val.(string)
		component.Ability = ability
	}

	return component
}
