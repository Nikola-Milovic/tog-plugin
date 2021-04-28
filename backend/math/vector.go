package math

import (
	"fmt"
	"math"
)

//Vector represents a Vector2 with X and Y coordinates
type Vector struct {
	X float32 `json:"x"` //update constants.Vectorx and y
	Y float32 `json:"y"`
}

//FromScalar ..
func FromScalar(v float32) Vector {
	return Vector{v, v}
}

func V(x, y float32) Vector {
	return Vector{X: x, Y: y}
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

//To returns vector from v to u. Same as u.Sub(v)
func (v Vector) To(u Vector) Vector {
	return Vector{X: u.X - v.X, Y: u.Y - v.Y}
}

func (v Vector) Len2() float32 {
	return v.X*v.X + v.Y*v.Y
}


//Lerp
func (v Vector) Lerp(b Vector, t float32) Vector {
	return v.To(b).MultiplyScalar(t).Add(v)
}

// Point converts Vec to Point
func (v Vector) Point() Point {
	return Point{int(v.X), int(v.Y)}
}

// Max uses MaxF on both components and returns resulting vector
func (v Vector) Max(u Vector) Vector {
	return Vector{
		MaxF(v.X, u.X),
		MaxF(v.Y, u.Y),
	}
}

// Min uses MinF on both components and returns resulting vector
func (v Vector) Min(u Vector) Vector {
	return Vector{
		MinF(v.X, u.X),
		MinF(v.Y, u.Y),
	}
}

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

//GetDistanceIncludingDiagonalVectors returns distance along with diagonals
func GetDistanceIncludingDiagonalVectors(c1 Vector, c2 Vector) float32 {

	r := math.Max(math.Abs(float64(c1.X-c2.X)), math.Abs(float64(c1.Y-c2.Y))) // TODO change to float32

	return float32(r)
}

func GetDistanceBetweenPoints(p1 Point, p2 Point) int {
	r := Max(AbsI(p1.X-p2.X), AbsI(p1.Y-p2.Y))

	return r
}