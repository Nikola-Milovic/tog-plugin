package components

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type PositionComponent struct {
	Position math.Vector
	Address  math.Point
	Radius   float32
}

func (m PositionComponent) ComponentName() string {
	return "PositionComponent"
}

func PositionComponentMaker(data interface{}, additionalData map[string]interface{}, world engine.WorldI) engine.Component {
	pos := additionalData["position"].(math.Vector)
	tag := additionalData["tag"].(int)
	start := additionalData["start"].(bool)

	posData := data.(map[string]interface{})

	radius := float32(posData["Radius"].(float64))

	if start && tag == 1 { // Used to place the other player at the other end of the screen
		pos.X = float32(800) - pos.X*32 //TODO change so it isn't 32 but rather constant X
	} else {
		pos.X = pos.X * 32 //Todo add constants
	}

	pos.Y = pos.Y * 32

	return PositionComponent{Position: pos, Radius: radius}
}
