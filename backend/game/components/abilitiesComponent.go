package components

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

type Ability map[string]interface{}

type AbilitiesComponent struct {
	Abilities map[string]Ability
}

func (a AbilitiesComponent) ComponentName() string {
	return "AbilitiesComponent"
}

func AbilitiesComponentMaker(data interface{}, abData interface{}, world *game.World) engine.Component {
	component := AbilitiesComponent{}

	//abilityDataMap := abData.(map[string]interface{})

	compData := data.([]interface{})

	component.Abilities = make(map[string]Ability, len(compData))

	for _, a := range compData {
		ab := a.(map[string]interface{})
		abilityID := ab["AbilityID"].(string)
		ability := make(map[string]interface{})

		component.Abilities[abilityID] = ability

		//If the ability should be available instantly, eg summons or buffs or something
		if _, ok := ab["InstantCast"]; ok {
			component.Abilities[abilityID]["last_activated"] = -10000
		}
	}

	return component
}

func onHitAbilityType(ability map[string]interface{}, id string, abData map[string]interface{}) {
	//		switch abilityDataMap[abilityID].(map[string]interface{})["Type"].(string) {
	// case "OnHit":
	// 	onHitAbilityType(ability, abilityID, abilityDataMap[abilityID].(map[string]interface{}), world*game.World)
	// }
}
