package engine

import (
	"fmt"
	"math"
)

//Vector represents a Vector2 with X and Y coordinates
type Vector struct {
	X float32 `json:"x"` //update constants.Vectorx and y
	Y float32 `json:"y"`
}

//New ..
func New(x, y float32) Vector {
	return Vector{x, y}
}

//FromScalar ..
func FromScalar(v float32) Vector {
	return Vector{v, v}
}

//Zero ..
func Zero() Vector {
	return Vector{0, 0}
}

//Unit ..
func Unit() Vector {
	return Vector{1, 1}
}

//Copy ..
func (v Vector) Copy() Vector {
	return Vector{v.X, v.Y}
}

//Add ..
func (v Vector) Add(v2 Vector) Vector {
	return Vector{v.X + v2.X, v.Y + v2.Y}
}

//Subtract ..
func (v Vector) Subtract(v2 Vector) Vector {
	return Vector{v.X - v2.X, v.Y - v2.Y}
}

//SubtractScalar ..
func (v Vector) SubtractScalar(val float32) Vector {
	return Vector{v.X - val, v.Y - val}
}

//Multiply ..
func (v Vector) Multiply(v2 Vector) Vector {
	return Vector{v.X * v2.X, v.Y * v2.Y}
}

//Divide ..
func (v Vector) Divide(v2 Vector) Vector {
	return Vector{v.X / v2.X, v.Y / v2.Y}
}

//MultiplyScalar ..
func (v Vector) MultiplyScalar(s float32) Vector {
	return Vector{v.X * s, v.Y * s}
}

//DivideScalar ..
func (v Vector) DivideScalar(s float32) Vector {
	return Vector{v.X / s, v.Y / s}
}

//Norm2 ..
func (v Vector) Norm2() float64 { return float64(v.Dot(v)) }

//Dot ..
func (v Vector) Dot(ov Vector) float32 { return v.X*ov.X + v.Y*ov.Y }

//Norm ..
func (v Vector) Norm() float32 { return float32(math.Sqrt(float64(v.Dot(v)))) }

//Magnitute ..
func (v Vector) Magnitute() float32 { return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))) }

// Normalize returns a unit vector in the same direction as v.
func (v Vector) Normalize() Vector {
	n2 := v.Norm2()
	if n2 == 0 {
		return Vector{0, 0}
	}
	return v.MultiplyScalar(float32(1 / math.Sqrt(n2)))
}

func (v Vector) String() string {
	return fmt.Sprintf("X : %v, Y :%v", v.X, v.Y)
}

// Distance returns the Euclidean distance between v and ov.
func (v Vector) Distance(ov Vector) float32 { return v.Subtract(ov).Norm() }

//GetDistance returns the distance the way unit would travel, without diagonals
func GetDistance(c1 Vector, c2 Vector) float32 {
	absX := c1.X - c2.X
	if absX < 0 {
		absX = -absX
	}
	absY := c1.Y - c2.Y
	if absY < 0 {
		absY = -absY
	}
	r := absX + absY

	return r
}

//GetDistanceIncludingDiagonal returns distance along with diagonals
func GetDistanceIncludingDiagonal(c1 Vector, c2 Vector) float32 {

	r := math.Max(math.Abs(float64(c1.X-c2.X)), math.Abs(float64(c1.Y-c2.Y)))

	return float32(r)
}
