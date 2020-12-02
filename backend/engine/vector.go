package engine

import (
	"fmt"
)

//Vector represents a Vector2 with X and Y coordinates
type Vector struct {
	x int `json:"pos_x"`
	y int `json:"pos_y"`
}

func (v Vector) String() string {
	return fmt.Sprintf("X : %v, Y :%v", v.x, v.y)
}
