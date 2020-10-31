package constants

import (
	"fmt"
)

//V2 represents a Vector2 with X and Y coordinates
type V2 struct {
	X int `json:"pos_x"`
	Y int `json:"pos_y"`
}

func (v V2) String() string {
	return fmt.Sprintf("X : %v, Y :%v", v.X, v.Y)
}
