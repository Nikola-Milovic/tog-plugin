package constants

import (
	"fmt"
	"math"
)

type V2 struct {
	X float64 `json:"pos_x"`
	Y float64 `json:"pos_y"`
}

func New(x, y float64) V2 {
	return V2{x, y}
}

func FromScalar(v float64) V2 {
	return V2{v, v}
}

func Zero() V2 {
	return V2{0, 0}
}

func Unit() V2 {
	return V2{1, 1}
}

func (v V2) Copy() V2 {
	return V2{v.X, v.Y}
}

func (v V2) Add(v2 V2) V2 {
	return V2{v.X + v2.X, v.Y + v2.Y}
}

func (v V2) Subtract(v2 V2) V2 {
	return V2{v.X - v2.X, v.Y - v2.Y}
}

func (v V2) Multiply(v2 V2) V2 {
	return V2{v.X * v2.X, v.Y * v2.Y}
}

func (v V2) Divide(v2 V2) V2 {
	return V2{v.X / v2.X, v.Y / v2.Y}
}

func (v V2) MultiplyScalar(s float64) V2 {
	return V2{v.X * s, v.Y * s}
}

func (v V2) DivideScalar(s float64) V2 {
	return V2{v.X / s, v.Y / s}
}

func (v V2) Norm2() float64 { return v.Dot(v) }

func (v V2) Dot(ov V2) float64 { return v.X*ov.X + v.Y*ov.Y }

func (v V2) Norm() float64 { return math.Sqrt(v.Dot(v)) }

// Normalize returns a unit vector in the same direction as v.
func (v V2) Normalize() V2 {
	n2 := v.Norm2()
	if n2 == 0 {
		return V2{0, 0}
	}
	return v.MultiplyScalar(1 / math.Sqrt(n2))
}

func (v V2) String() string {
	return fmt.Sprintf("X : %v, Y :%v", v.X, v.Y)
}

// Distance returns the Euclidean distance between v and ov.
func (v V2) Distance(ov V2) float64 { return v.Subtract(ov).Norm() }

// func (v *V2) Add(v2 V2) {
// 	v.X += v2.X
// 	v.Y += v2.Y
// }

// func (v *V2) Subtract(v2 V2) {
// 	v.X -= v2.X
// 	v.Y -= v2.Y
// }

// func (v *V2) Multiply(v2 V2) {
// 	v.X *= v2.X
// 	v.Y *= v2.Y
// }

// func (v *V2) Divide(v2 V2) {
// 	v.X /= v2.X
// 	v.Y /= v2.Y
// }

// func (v *V2) MultiplyScalar(s float64) {
// 	v.X *= s
// 	v.Y *= s
// }

// func (v *V2) DivideScalar(s float64) {
// 	v.X /= s
// 	 v.Y /= s
// }

// func (v *V2) Norm2() float64 { return v.Dot(v) }

// func (v *V2) Dot(ov V2) float64 { return v.X*ov.X + v.Y*ov.Y }

// func (v *V2) Norm() float64 { return math.Sqrt(v.Dot(v)) }
