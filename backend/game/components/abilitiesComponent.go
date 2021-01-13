package components

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type Ability map[string]interface{}

type AbilitiesComponent struct {
	Abilities map[string]Ability
}

func (a AbilitiesComponent) ComponentName() string {
	return "AbilitiesComponent"
}

func AbilitiesComponentMaker(data interface{}, abilityData interface{}) engine.Component {
	component := AbilitiesComponent{}
	
	compData := data.([]interface{})

	component.Abilities = make(map[string]Ability, len(compData))

	for _, a := range compData {
		ab := a.(map[string]interface{})
		ability := make(map[string]interface{})
		abilityID := ab["AbilityID"].(string)
		component.Abilities[abilityID] = ability

		//If the ability should be available instantly, eg summons or buffs or something
		if _, ok := ab["InstantCast"]; ok {
			component.Abilities[abilityID]["last_activated"] = -10000
		}
	}

	return component
}
