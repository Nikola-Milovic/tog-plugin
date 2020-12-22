package components

import "github.com/Nikola-Milovic/tog-plugin/engine"

type PositionComponent struct {
	Position engine.Vector
}

func (m PositionComponent) ComponentName() string {
	return "PositionComponent"
}

func PositionComponentMaker(data interface{}) engine.Component {
	return PositionComponent{}
}
