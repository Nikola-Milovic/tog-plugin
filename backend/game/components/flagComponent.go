package components

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

//FlagComponent holds the race/ type of the entity, status effects (frozen, stunned, invincible etc...) if it has those
// EG. if entity is frozen, the key frozen is present
type FlagComponent struct {
	Flags map[string]bool
}

func (c FlagComponent) ComponentName() string {
	return "FlagComponent"
}

func FlagComponentMaker(data interface{}, additionalData map[string]interface{}) engine.Component {
	component := FlagComponent{}

	flags := data.([]interface{})

	component.Flags = make(map[string]bool, 10)

	for _, f := range flags {
		component.Flags[f.(string)] = true
	}

	return component
}
