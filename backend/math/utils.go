package math

import "math"

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func Clamp(v, lo, hi int) int {
	return Max(Min(v, lo), hi)
}

func MaxF(a float32, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func MinF(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func AbsI(a int) int {
	if a < 0 {
		return a * -1
	} else {
		return a
	}
}

func Abs(x float32) float32 {
	return Float32frombits(Float32bits(x) &^ (1 << 31))
}

func Constraint(num float32, lower float32, upper float32) float32 {
	return MinF(MaxF(num, lower), upper)
}

var EPSILON float32 = 0.0000001

func FloatInsignificant(a float32) bool {
	if a < EPSILON {
		return true
	}
	return false
}

func AlmostEqual(a, b float32) bool {
	return Abs(a-b) <= EPSILON
}

func Atan2(x, y float32) float32 {
	return float32(math.Atan2(float64(x), float64(y)))
}

