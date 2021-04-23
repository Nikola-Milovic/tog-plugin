package components

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type PositionComponent struct {
	Position    math.Vector
	BoundingBox math.Vector
}

func (m PositionComponent) ComponentName() string {
	return "PositionComponent"
}

func PositionComponentMaker(data interface{}, additionalData map[string]interface{}, world engine.WorldI) engine.Component {
	pos := additionalData["position"].(math.Vector)
	tag := additionalData["tag"].(int)
	start := additionalData["start"].(bool)

	posData := data.(map[string]interface{})

	bbox := posData["BoundingBox"].(map[string]interface{})
	boundingBox := math.Vector{X: float32(bbox["x"].(float64)), Y: float32(bbox["y"].(float64))}

	if start && tag == 1 { // Used to place the other player at the other end of the screen
		fmt.Printf("X : %d, Y : %d \n", pos.X, pos.Y)
		pos.X = float32(800) - pos.X*32 //TODO change so it isn't 32 but rather constant X
		fmt.Printf("new X : %d, Y : %d \n", pos.X, pos.Y)
	} else {
		pos.X = pos.X * 32 //Todo add constants
	}

	pos.Y = pos.Y * 32

	return PositionComponent{Position: pos, BoundingBox: boundingBox}
}
