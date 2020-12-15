package engine

import (
	"fmt"
)

//Vector represents a Vector2 with X and Y coordinates
type Vector struct {
	X int `json:"x"` //update constants.Vectorx and y
	Y int `json:"y"`
}

func (v Vector) String() string {
	return fmt.Sprintf("X : %v, Y :%v", v.X, v.Y)
}
