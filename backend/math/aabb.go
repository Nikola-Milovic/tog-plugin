package math

import (
	"fmt"
	"image"
	"math"
)

// AABB is a  rectangle aligned with the axes of the coordinate system. It is defined by two
// points, Min and Max.
//
// The invariant should hold, that Max's components are greater or equal than Min's components
// respectively.
type AABB struct {
	Min, Max Vector
}

// ZA is zero value AABB
var ZA AABB

// A returns a new AABB with given the Min and Max coordinates.
//
// Note that the returned rectangle is not automatically normalized.
func A(minX, minY, maxX, maxY float32) AABB {
	return AABB{
		Min: Vector{minX, minY},
		Max: Vector{maxX, maxY},
	}
}

// Square returns AABB with center in c and width and height both equal to size * 2
func Square(c Vector, size float32) AABB {
	return AABB{Vector{c.X - size, c.Y - size}, Vector{c.X + size, c.Y + size}}
}

// Centered returns rect witch center is equal to c, width equal to w, likewise height equal to h
func Centered(c Vector, w, h float32) AABB {
	w, h = w/2, h/2
	return AABB{Vector{X: c.X - w, Y: c.Y - h}, Vector{X: c.X + w, Y: c.Y + h}}
}

// FromRect converts image.Rectangle to AABB
func FromRect(a image.Rectangle) AABB {
	return AABB{
		Vector{float32(a.Min.X), float32(a.Min.Y)},
		Vector{float32(a.Max.X), float32(a.Max.Y)},
	}
}

// ToImage converts AABB to image.AABB
func (a AABB) ToImage() image.Rectangle {
	return image.Rect(
		int(a.Min.X),
		int(a.Min.Y),
		int(a.Max.X),
		int(a.Max.Y),
	)
}

// ToVector converts AABB to vec where x is AABB width anc y is rect Height
func (a AABB) ToVector() Vector {
	return a.Min.To(a.Max)
}

// String returns the string representation of the AABB.
func (a AABB) String() string {
	return fmt.Sprintf("A(%v %v %v %v)", ff(a.Min.X), ff(a.Min.Y), ff(a.Max.X), ff(a.Max.Y))
}

// Norm returns the AABB in normal form, such that Max is component-wise greater or equal than Min.
func (a AABB) Norm() AABB {
	return AABB{
		Min: Vector{
			MinF(a.Min.X, a.Max.X),
			MinF(a.Min.Y, a.Max.Y),
		},
		Max: Vector{
			MaxF(a.Min.X, a.Max.X),
			MaxF(a.Min.Y, a.Max.Y),
		},
	}
}

// W returns the width of the AABB.
func (a AABB) W() float32 {
	return a.Max.X - a.Min.X
}

// H returns the height of the AABB.
func (a AABB) H() float32 {
	return a.Max.Y - a.Min.Y
}

// Size returns the vector of width and height of the AABB.
func (a AABB) Size() Vector {
	return Vector{a.W(), a.H()}
}

// Area returns the area of a. If a is not normalized, area may be negative.
func (a AABB) Area() float32 {
	return a.W() * a.H()
}

// Center returns the position of the center of the AABB.
func (a AABB) Center() Vector {
	return a.Min.Lerp(a.Max, 0.5)
}

// Moved returns the AABB moved (both Min and Max) by the given vector delta.
func (a AABB) Moved(delta Vector) AABB {
	return AABB{
		Min: a.Min.Add(delta),
		Max: a.Max.Add(delta),
	}
}

// Resized returns the AABB resized to the given size while keeping the position of the given
// anchor.
//
//   a.Resized(a.Min, size)      // resizes while keeping the position of the lower-left corner
//   a.Resized(a.Max, size)      // same with the top-right corner
//   a.Resized(a.Center(), size) // resizes around the center
func (a AABB) Resized(anchor, size Vector) AABB {
	fraction := Vector{size.X / a.W(), size.Y / a.H()}
	return AABB{
		Min: anchor.Add(a.Min.Subtract(anchor).Multiply(fraction)),
		Max: anchor.Add(a.Max.Subtract(anchor).Multiply(fraction)),
	}
}

// ResizedMin returns the AABB resized to the given size while keeping the position of the AABB's
// Min.
//
// Sizes of zero area are safe here.
func (a AABB) ResizedMin(size Vector) AABB {
	return AABB{
		Min: a.Min,
		Max: a.Min.Add(size),
	}
}

// Contains checks whether a vector u is contained within this AABB (including it's borders).
func (a AABB) Contains(u Vector) bool {
	return a.Min.X <= u.X && u.X <= a.Max.X && a.Min.Y <= u.Y && u.Y <= a.Max.Y
}

// Union returns the minimal AABB which covers both a and s. AABBs a and s must be normalized.
func (a AABB) Union(s AABB) AABB {
	return A(
		MinF(a.Min.X, s.Min.X),
		MinF(a.Min.Y, s.Min.Y),
		MaxF(a.Max.X, s.Max.X),
		MaxF(a.Max.Y, s.Max.Y),
	)
}

// Intersect returns the maximal AABB which is covered by both a and s. AABBs a and s must be normalized.
//
// If a and s don't overlap, this function returns a zero-rectangle.
func (a AABB) Intersect(s AABB) AABB {
	t := A(
		MaxF(a.Min.X, s.Min.X),
		MaxF(a.Min.Y, s.Min.Y),
		MinF(a.Max.X, s.Max.X),
		MinF(a.Max.Y, s.Max.Y),
	)

	if t.Min.X >= t.Max.X || t.Min.Y >= t.Max.Y {
		return AABB{}
	}

	return t
}

// Intersects returns whether or not the given AABB intersects at any point with this AABB.
//
// This function is overall about 5x faster than Intersect, so it is better
// to use if you have no need for the returned AABB from Intersect.
func (a AABB) Intersects(s AABB) bool {
	return !(s.Max.X < a.Min.X ||
		s.Min.X > a.Max.X ||
		s.Max.Y < a.Min.Y ||
		s.Min.Y > a.Max.Y)
}

// Vertices returns a slice of the four corners which make up the rectangle.
func (a AABB) Vertices() [4]Vector {
	return [4]Vector{
		a.Min,
		{a.Min.X, a.Max.Y},
		a.Max,
		{a.Max.X, a.Min.Y},
	}
}

// LocalVertices creates array of vertices relative to center of rect
func (a AABB) LocalVertices() [4]Vector {
	v := a.Vertices()
	c := a.Center()

	for i, e := range v {
		v[i] = e.Subtract(c)
	}

	return v
}

// VectorBounds gets the smallest rectangle in witch all provided points fit in
func VectorBounds(vectors ...Vector) (base AABB) {
	base.Min.X = math.MaxFloat32
	base.Min.Y = math.MaxFloat32
	base.Max.X = -math.MaxFloat32
	base.Max.Y = -math.MaxFloat32

	for _, v := range vectors {
		if base.Min.X > v.X {
			base.Min.X = v.X
		}
		if base.Min.Y > v.Y {
			base.Min.Y = v.Y
		}
		if base.Max.X < v.X {
			base.Max.X = v.X
		}
		if base.Max.Y < v.Y {
			base.Max.Y = v.Y
		}
	}

	return base
}

// Clamp clamps Vector inside AABB area
func (a AABB) Clamp(v Vector) Vector {
	return Vector{
		MaxF(MinF(v.X, a.Max.X), a.Min.X),
		MaxF(MinF(v.Y, a.Max.Y), a.Min.Y),
	}
}

// Flatten returns AABB flattened into Array, values are
// in same order as they would be stored on stack
func (a AABB) Flatten() [4]float32 {
	return [...]float32{a.Min.X, a.Min.Y, a.Max.X, a.Max.Y}
}

// Mutator is similar to Iterator but this gives option to mutate
// state of AABB trough Array Entries
func (a *AABB) Mutator() [4]*float32 {
	return [...]*float32{&a.Min.X, &a.Min.Y, &a.Max.X, &a.Max.Y}
}

// Deco returns edge values
func (a *AABB) Deco() (left, bottom, right, top float32) {
	return a.Min.X, a.Min.Y, a.Max.X, a.Max.Y
}

// Fits reports whether a fits into b so that a.Intersect(b) == a
func (a AABB) Fits(b AABB) bool {
	return b.Min.X <= a.Min.X && b.Min.Y <= a.Min.Y && b.Max.X >= a.Max.X && b.Max.Y >= a.Max.Y
}

//// IntersectsAABB returns whether circle intersects AABB
//func (a AABB) IntersectCircle(c Circ) bool {
//	r, l, t, b := c.C.X > a.Min.X, c.C.X < a.Max.X, c.C.Y > a.Min.Y, c.C.Y < a.Max.Y
//
//	if t && b {
//		if r {
//			return l || c.C.X-c.R <= a.Max.X
//		} else {
//			return l && c.C.X+c.R >= a.Min.X
//		}
//	}
//
//	if l && r {
//		return (t && c.C.Y-c.R <= a.Max.Y) || b && c.C.Y+c.R >= a.Min.Y
//	}
//
//	v := a.Vertices()
//	r2 := c.R * c.R
//
//	if b {
//		if (l && c.C.To(v[0]).Len2() <= r2) || c.C.To(v[3]).Len2() <= r2 {
//			return true
//		}
//	} else {
//		if (l && c.C.To(v[1]).Len2() <= r2) || c.C.To(v[2]).Len2() <= r2 {
//			return true
//		}
//	}
//
//	return false
//}

func ff(f float32) string {
	return fmt.Sprintf("%.2f", f)
}
