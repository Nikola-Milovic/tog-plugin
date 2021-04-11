package engine

// Component interface should be inherited by all components
type Component interface {
	ComponentName() string
}

// ComponentMaker represents a type of function that creates a specific Component
type ComponentMakerFun func(interface{}, map[string]interface{}, WorldI) Component

type ComponentMaker interface {
	RegisterComponentMaker(string, ComponentMakerFun)
	AddComponents(map[string]interface{}, string, map[string]interface{})
}
