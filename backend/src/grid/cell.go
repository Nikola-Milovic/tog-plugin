package grid

import (
	"github.com/Nikola-Milovic/tog-plugin/src/astar"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

type Cell struct {
	Position   constants.V2
	isOccupied bool
	grid       *Grid
}

func (c *Cell) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		if n := c.grid.CellAt(c.Position.X+offset[0], c.Position.Y+offset[1]); n != nil &&
			!n.isOccupied {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (c *Cell) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

func (c *Cell) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*Cell)
	absX := toT.Position.X - c.Position.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Position.Y - c.Position.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}
