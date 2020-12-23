package components

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type EffectsComponent struct {
	Effects []map[string]interface{}
}

func (e EffectsComponent) ComponentName() string {
	return "EffectsComponent"
}

func EffectsComponentMaker(data interface{}) engine.Component {
	component := EffectsComponent{}

	component.Effects = make([]map[string]interface{}, 0, 10)

	return component
}
