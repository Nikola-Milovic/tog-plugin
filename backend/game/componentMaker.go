package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
)

// ComponentMaker represents a type of function that creates a specific Component
type ComponentMaker func(interface{}) engine.Component

// UniqueComponentMaker represents a type of function that creates a specific UniqueComponent
type UniqueComponentMaker func(interface{}, interface{}) engine.Component

type ComponentManager struct {
	World                   *World
	ComponentRegistry       map[string]ComponentMaker
	UniqueComponentRegistry map[string]UniqueComponentMaker
}

func (cm *ComponentManager) AddComponents(data map[string]interface{}, id string) {
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
			component := maker(data, cm.World.AbilityDataMap)
			cm.World.EntityManager.ObjectPool.AddUniqueComponent(component, id)
			continue
		}
		maker, ok := cm.ComponentRegistry[key] // Move.EntityManagerentComponentMaker, returns a Move.EntityManagerentComponent
		if !ok {
			panic(fmt.Sprintf("No registered maker for the component %s", key))
		}
		component := maker(data)
		cm.World.ObjectPool.AddComponent(component)

	}
}

func (cm *ComponentManager) RegisterComponentMaker(componentName string, maker ComponentMaker) {
	if _, ok := cm.ComponentRegistry[componentName]; ok {
		panic(fmt.Sprintf("Component maker for component %v is already registered", componentName))
	}
	cm.ComponentRegistry[componentName] = maker
}

func (cm *ComponentManager) RegisterUniqueComponentMaker(componentName string, maker UniqueComponentMaker) {
	if _, ok := cm.UniqueComponentRegistry[componentName]; ok {
		panic(fmt.Sprintf("Unique component maker for component %v is already registered", componentName))
	}
	cm.UniqueComponentRegistry[componentName] = maker
}

func CreateComponentManager(w *World) *ComponentManager {
	cm := ComponentManager{}

	cm.World = w
	cm.ComponentRegistry = make(map[string]ComponentMaker)
	cm.UniqueComponentRegistry = make(map[string]UniqueComponentMaker, 2)

	return &cm
}
