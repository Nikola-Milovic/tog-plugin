package components

import "github.com/Nikola-Milovic/tog-plugin/engine"

type PositionComponent struct {
	Position engine.Vector
}

func (m PositionComponent) ComponentName() string {
	return "PositionComponent"
}

func PositionComponentMaker(data interface{}, additionalData map[string]interface{}, world engine.WorldI) engine.Component {
	pos := additionalData["position"].(engine.Vector)
	tag := additionalData["tag"].(int)
	start := additionalData["start"].(bool)

	if start && tag == 1 { // Used to place the other player at the other end of the screen
		pos.X = 800/32 - pos.X //Todo add constants
	}

	return PositionComponent{Position: pos}
}
