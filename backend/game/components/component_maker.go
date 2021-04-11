package components

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type ComponentMaker struct {
	World             engine.WorldI
	ComponentRegistry map[string]engine.ComponentMakerFun
	CommonComponents  []string
}

func (cm *ComponentMaker) AddComponents(componentData map[string]interface{}, id string, additionalData map[string]interface{}) {
	components, ok := componentData["Components"].(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("Unit has no components %s", id))
	}

	//Eg key = Move.EntityManagerentComponent, data is Move.EntityManagerentSpeed, Move.EntityManagerentType etc
	for key, data := range components {
		// if key == "AbilitiesComponent" {
		// 	maker, ok := cm.UniqueComponentRegistry[key]
		// 	if !ok {
		// 		panic(fmt.Sprintf("No registered maker for the component %s for index %v", key))
		// 	}
		// 	component := maker(data, cm.World.GetAbilityDataMap(), cm.World, id)
		// 	cm.World.GetObjectPool().AddUniqueComponent(component, id)
		// 	continue
		// }
		maker, ok := cm.ComponentRegistry[key] // Move.EntityManagerentComponentMaker, returns a Move.EntityManagerentComponent
		if !ok {
			panic(fmt.Sprintf("No registered maker for the component %s", key))
		}
		component := maker(data, additionalData, cm.World)
		cm.World.GetObjectPool().AddComponent(component)

	}

	for _, comp := range cm.CommonComponents {
		maker, ok := cm.ComponentRegistry[comp]
		if !ok {
			panic(fmt.Sprintf("No registered maker for the component %s", comp))
		}
		component := maker(componentData, additionalData, cm.World)
		cm.World.GetObjectPool().AddComponent(component)
	}

}

func (cm *ComponentMaker) RegisterComponentMaker(componentName string, maker engine.ComponentMakerFun) {
	if _, ok := cm.ComponentRegistry[componentName]; ok {
		panic(fmt.Sprintf("Component maker for component %v is already registered", componentName))
	}
	cm.ComponentRegistry[componentName] = maker
}

func CreateComponentMaker(w engine.WorldI) *ComponentMaker {
	cm := ComponentMaker{}

	cm.World = w
	cm.ComponentRegistry = make(map[string]engine.ComponentMakerFun)
	cm.CommonComponents = []string{"EffectsComponent"}

	return &cm
}
