package constants

import (
	"fmt"
	"math"
)

//V2 represents a Vector2 with X and Y coordinates
type V2 struct {
	X float32 `json:"pos_x"`
	Y float32 `json:"pos_y"`
}

//New ..
func New(x, y float32) V2 {
	return V2{x, y}
}

//FromScalar ..
func FromScalar(v float32) V2 {
	return V2{v, v}
}

//Zero ..
func Zero() V2 {
	return V2{0, 0}
}

//Unit ..
func Unit() V2 {
	return V2{1, 1}
}

//Copy ..
func (v V2) Copy() V2 {
	return V2{v.X, v.Y}
}

//Add ..
func (v V2) Add(v2 V2) V2 {
	return V2{v.X + v2.X, v.Y + v2.Y}
}

//Subtract ..
func (v V2) Subtract(v2 V2) V2 {
	return V2{v.X - v2.X, v.Y - v2.Y}
}

//SubtractScalar ..
func (v V2) SubtractScalar(val float32) V2 {
	return V2{v.X - val, v.Y - val}
}

//Multiply ..
func (v V2) Multiply(v2 V2) V2 {
	return V2{v.X * v2.X, v.Y * v2.Y}
}

//Divide ..
func (v V2) Divide(v2 V2) V2 {
	return V2{v.X / v2.X, v.Y / v2.Y}
}

//MultiplyScalar ..
func (v V2) MultiplyScalar(s float32) V2 {
	return V2{v.X * s, v.Y * s}
}

//DivideScalar ..
func (v V2) DivideScalar(s float32) V2 {
	return V2{v.X / s, v.Y / s}
}

//Norm2 ..
func (v V2) Norm2() float64 { return float64(v.Dot(v)) }

//Dot ..
func (v V2) Dot(ov V2) float32 { return v.X*ov.X + v.Y*ov.Y }

//Norm ..
func (v V2) Norm() float32 { return float32(math.Sqrt(float64(v.Dot(v)))) }

//Magnitute ..
func (v V2) Magnitute() float32 { return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))) }

// Normalize returns a unit vector in the same direction as v.
func (v V2) Normalize() V2 {
	n2 := v.Norm2()
	if n2 == 0 {
		return V2{0, 0}
	}
	return v.MultiplyScalar(float32(1 / math.Sqrt(n2)))
}

func (v V2) String() string {
	return fmt.Sprintf("X : %v, Y :%v", v.X, v.Y)
}

// Distance returns the Euclidean distance between v and ov.
func (v V2) Distance(ov V2) float32 { return v.Subtract(ov).Norm() }
