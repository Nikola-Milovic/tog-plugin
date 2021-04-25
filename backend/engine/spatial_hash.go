package engine

import (
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type SpatialHash struct {
	nodeSize math.Vector
	w, h     int
	Nodes    []SpatialNode
}

// CreateSpatialHash is SpatialHash constructor
func CreateSpatialHash(w, h int, tileSize math.Vector) *SpatialHash {
	return &SpatialHash{
		nodeSize: math.V(1, 1).Divide(tileSize),
		w:        w,
		h:        h,
		Nodes:    make([]SpatialNode, h*w),
	}
}

func (h *SpatialHash) Insert(adr math.Point, pos math.Vector, id, group int) math.Point {
	adr = h.Adr(pos)
	h.Nodes[getIndex(adr.X, adr.Y, h.w)].Insert(id, group)
	return adr
}

// Remove removes shape from SpatialHash. If operation fails, false is returned
func (h *SpatialHash) Remove(adr math.Point, id, group int) bool {
	return h.Nodes[getIndex(adr.X, adr.Y, h.w)].Remove(id, group)
}

// Update updates state of object if it changed quadrant, if operation fails, false is returned
func (h *SpatialHash) Update(old math.Point, pos math.Vector, id, group int) math.Point {
	p := h.Adr(pos)
	if old == p {
		return p
	}

	if h.Nodes[getIndex(old.X, old.Y, h.w)].Remove(id, group) {
		h.Nodes[getIndex(p.X, p.Y, h.w)].Insert(id, group)
		return p
	}

	return p
}

// Query returns colliding shapes with given rect
func (h *SpatialHash) Query(rect math.AABB, coll []int, group int, including bool) []int {
	max := h.Adr(rect.Max).Add(math.P(2, 2)).Min(math.P(h.w, h.h))
	min := h.Adr(rect.Min).Sub(math.P(1, 1)).Max(math.P(0, 0))

	for y := min.Y; y < max.Y; y++ {
		for x := min.X; x < max.X; x++ {
			n := &h.Nodes[getIndex(x, y, h.w)]
			if n.Count != 0 {
				if group == -1 {
					coll = n.CollectAll(coll)
				} else {
					coll = n.Collect(group, including, coll)
				}
			}
		}
	}

	return coll
}

func getIndex(x, y, stride int) int {
	return y*stride + x
}

// Adr returns node, position belongs to
func (h *SpatialHash) Adr(pos math.Vector) math.Point {
	// we want this inlined
	x, y := int(pos.X*h.nodeSize.X), int(pos.Y*h.nodeSize.Y)
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= h.w {
		x = h.w - 1
	}
	if y >= h.h {
		y = h.h - 1
	}
	return math.P(x, y)
}
