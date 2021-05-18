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

func (v Vector) Crossed() Vector {
	return Vector{X: -v.Y, Y: v.X}
}

//Multiply ..
func (v Vector) Multiply(v2 Vector) Vector {
	return Vector{v.X * v2.X, v.Y * v2.Y}
}

//Divide ..
func (v Vector) Divide(v2 Vector) Vector {
	return Vector{v.X / v2.X, v.Y / v2.Y}
}

func (v Vector) Trunctate(max float32) Vector {
	i := max / v.Magnitute()
	if i < 1.0 {
		i = 1.0
	}
	return v.MultiplyScalar(i)
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
func (v Vector) Norm2() float32 { return v.Dot(v) }

//Dot ..
func (v Vector) Dot(ov Vector) float32 { return v.X*ov.X + v.Y*ov.Y }

func (v Vector) AngleTo(ov Vector) float32 {
	angle := Atan2(ov.Y, ov.X) - Atan2(v.Y, v.X)
	if angle > 2*Pi {
		angle -= 2 * Pi
	} else if angle <= -2*Pi {
		angle += 2 * Pi
	}
	return angle

	//v1 := v.Normalize()
	//v2 := ov.Normalize()
	//
	//return v1.X*v2.X + v1.Y*v2.Y

	//	angle := v.Dot(ov) / (v.Magnitute() * ov.Magnitute())
	//	// prevent NaN
	//	if angle > 1. {
	//		angle = angle - 2
	//	} else if angle < -1. {
	//		angle = angle + 2
	//	}
	//	return angle
}

// Crossed returns the cross product of two vectors.
func (a Vector) Cross( b Vector) Vector {
	return Vector{
		a.Y*b.X - a.X*b.Y,
		a.X*b.Y - a.Y*b.X,
	}
}

// ApproxEqual reports whether v and ov are equal within a small epsilon.
func (v Vector) ApproxEqual(ov Vector) bool {
	const epsilon = 1e-16
	return Abs(v.X-ov.X) < epsilon && Abs(v.Y-ov.Y) < epsilon
}

func (v Vector) PerpendicularClockwise() Vector {
	return Vector{v.Y, -v.X}
}

func (v Vector) PerpendicularCounterClockwise() Vector {
	return Vector{-v.Y, v.X}
}

//Norm ..
func (v Vector) Norm() float32 { return float32(math.Sqrt(float64(v.Dot(v)))) }

//Magnitute ..
func (v Vector) Magnitute() float32 { return Sqrt(v.X*v.X + v.Y*v.Y) }

// Normalize returns a unit vector in the same direction as v.
func (v Vector) Normalize() Vector {
	n2 := v.Norm2()
	if n2 == 0 {
		return Vector{0, 0}
	}
	return v.MultiplyScalar(1 / Sqrt(n2))
}

func (v Vector) String() string {
	return fmt.Sprintf("X: %.2f, Y:%.2f", v.X, v.Y)
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
