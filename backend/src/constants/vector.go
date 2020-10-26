package constants

import "fmt"

type V2 struct {
	X float32
	Y float32
}

func New(x, y float32) V2 {
	return V2{x, y}
}

func FromScalar(v float32) V2 {
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

func (v V2) MultiplyScalar(s float32) V2 {
	return V2{v.X * s, v.Y * s}
}

func (v V2) DivideScalar(s float32) V2 {
	return V2{v.X / s, v.Y / s}
}

func (v V2) String() string {
	return fmt.Sprintf("%v:%v", v.X, v.Y)
}
