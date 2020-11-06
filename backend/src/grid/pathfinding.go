package grid

import (
	"container/heap"

	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

// astar is an A* pathfinding implementation.

// node is a wrapper to store A* data for a Cell node.
type node struct {
	Cell   Cell
	cost   int
	rank   int
	parent *node
	open   bool
	closed bool
	index  int
}

// nodeMap is a collection of nodes keyed by Cell nodes for quick reference.
type nodeMap map[Cell]*node

// get gets the Cell object wrapped in a node, instantiating if required.
func (nm nodeMap) get(p Cell) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{
			Cell: p,
		}
		nm[p] = n
	}
	return n
}

// Path calculates a short path and the distance between the two Cell nodes.
//
// If no path is found, found will be false.
func Path(from, to Cell) (path []constants.V2, distance int, found bool) {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	fromNode := nm.get(from)
	fromNode.open = true
	heap.Push(nq, fromNode)
	for {
		if nq.Len() == 0 {
			// There's no path, return found false.
			return
		}
		current := heap.Pop(nq).(*node)
		current.open = false
		current.closed = true

		if current == nm.get(to) {
			// Found a path to the goal.
			p := []constants.V2{}
			curr := current
			for curr != nil {
				p = append(p, curr.Cell.Position)
				curr = curr.parent
			}
			return p[:len(p)-1], current.cost, true
		}

		for _, neighbor := range current.Cell.PathNeighbors() {
			cost := current.cost + 1
			neighborNode := nm.get(neighbor)

			if !neighborNode.open && !neighborNode.closed {
				neighborNode.cost = cost
				neighborNode.open = true
				neighborNode.rank = cost + neighbor.PathEstimatedCost(to)
				neighborNode.parent = current
				heap.Push(nq, neighborNode)
			}
		}
	}
}
