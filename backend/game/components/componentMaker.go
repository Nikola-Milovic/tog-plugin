package components

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

type ComponentMaker struct {
	World                   *game.World
	ComponentRegistry       map[string]engine.ComponentMakerFun
	UniqueComponentRegistry map[string]engine.UniqueComponentMakerFun
}

func (cm *ComponentMaker) AddComponents(data map[string]interface{}, id string, additionalData map[string]interface{}) {
	components, ok := data["Components"].(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("Unit has no components %s", id))
	}

	//Eg key = Move.EntityManagerentComponent, data is Move.EntityManagerentSpeed, Move.EntityManagerentType etc
	for key, data := range components {
		if key == "AbilitiesComponent" {
			maker, ok := cm.UniqueComponentRegistry[key]
			if !ok {
				panic(fmt.Sprintf("No registered maker for the component %s for index %v", key))
			}
			component := maker(data, cm.World.AbilityDataMap, cm.World)
			cm.World.EntityManager.ObjectPool.AddUniqueComponent(component, id)
			continue
		}
		maker, ok := cm.ComponentRegistry[key] // Move.EntityManagerentComponentMaker, returns a Move.EntityManagerentComponent
		if !ok {
			panic(fmt.Sprintf("No registered maker for the component %s", key))
		}
		component := maker(data, additionalData)
		cm.World.ObjectPool.AddComponent(component)

	}
}

func (cm *ComponentMaker) RegisterComponentMaker(componentName string, maker engine.ComponentMakerFun) {
	if _, ok := cm.ComponentRegistry[componentName]; ok {
		panic(fmt.Sprintf("Component maker for component %v is already registered", componentName))
	}
	cm.ComponentRegistry[componentName] = maker
}

func (cm *ComponentMaker) RegisterUniqueComponentMaker(componentName string, maker engine.UniqueComponentMakerFun) {
	if _, ok := cm.UniqueComponentRegistry[componentName]; ok {
		panic(fmt.Sprintf("Unique component maker for component %v is already registered", componentName))
	}
	cm.UniqueComponentRegistry[componentName] = maker
}

func CreateComponentMaker(w *game.World) *ComponentMaker {
	cm := ComponentMaker{}

	cm.World = w
	cm.ComponentRegistry = make(map[string]engine.ComponentMakerFun)
	cm.UniqueComponentRegistry = make(map[string]engine.UniqueComponentMakerFun, 2)

	return &cm
}
