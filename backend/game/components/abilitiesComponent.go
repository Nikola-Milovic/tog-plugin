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

func AbilitiesComponentMaker(data interface{}, abData interface{}, world engine.WorldI, id string) engine.Component {
	component := AbilitiesComponent{}

	//abilityDataMap := abData.(map[string]interface{})

	compData := data.([]interface{})

	component.Abilities = make(map[string]Ability, len(compData))

	abilityDataMap := world.GetAbilityDataMap()

	for _, a := range compData {
		ab := a.(map[string]interface{})
		abilityID := ab["AbilityID"].(string)
		ability := make(map[string]interface{})

		switch abilityDataMap[abilityID]["Type"].(string) {
		case "onHit":
			onHitAbilityType(ability, id, abilityDataMap[abilityID], world)
		default:
			//If the ability should be available instantly, eg summons or buffs or something
			if _, ok := ab["InstantCast"]; ok {
				ability["last_activated"] = -10000
			}
		}

		component.Abilities[abilityID] = ability
	}

	return component
}

func onHitAbilityType(ability map[string]interface{}, entityID string, abData map[string]interface{}, world engine.WorldI) {
	op := world.GetObjectPool()

	index := world.GetEntityManager().IndexMap[entityID]

	atkComp := op.Components["AttackComponent"][index].(AttackComponent)
	atkComp.OnHit = abData["Effect"].(string)
	op.Components["AttackComponent"][index] = atkComp
}
