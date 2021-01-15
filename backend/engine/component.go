package engine

// Component interface should be inherited by all components
type Component interface {
	ComponentName() string
}

// ComponentMaker represents a type of function that creates a specific Component
type ComponentMakerFun func(interface{}, map[string]interface{}) Component

// UniqueComponentMaker represents a type of function that creates a specific UniqueComponent
type UniqueComponentMakerFun func(interface{}, interface{}, *World) Component

type ComponentMaker interface {
	RegisterComponentMaker(string, ComponentMakerFun)
	RegisterUniqueComponentMaker(string, UniqueComponentMakerFun)
}
