package grid

import (
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type Tile struct {
	Position      math.Vector
	isOccupied    bool
	isGoal        bool // is already a goal of another unit
	occupiedIndex int
}

type FlowTile struct {
	Direction math.Vector
}

//instead of direct neightbours https://gamedevelopment.tutsplus.com/tutorials/how-to-speed-up-a-pathfinding-with-the-jump-point-search-algorithm--gamedev-5818
