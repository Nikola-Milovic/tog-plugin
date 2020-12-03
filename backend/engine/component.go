package engine

// Component interface should be inherited by all components
type Component interface {
	ComponentName() string
}

type ComponentMaker func(interface{}) Component
