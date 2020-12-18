package engine

// Component interface should be inherited by all components
type Component interface {
	ComponentName() string
}

// ComponentMaker represents a type of function that creates a specific Component
type ComponentMaker func(interface{}) Component
