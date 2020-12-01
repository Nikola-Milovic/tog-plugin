package engine

type Cell struct {
	Position   V2
	isOccupied bool
	grid       *Grid
	Index      int
}

//instead of direct neightbours https://gamedevelopment.tutsplus.com/tutorials/how-to-speed-up-a-pathfinding-with-the-jump-point-search-algorithm--gamedev-5818

func (c Cell) PathNeighbors() []Cell {
	return c.grid.GetNeighbours(c.Position.X, c.Position.Y)
}

func (c Cell) PathEstimatedCost(co Cell) int {
	absX := co.Position.X - c.Position.X
	if absX < 0 {
		absX = -absX
	}
	absY := co.Position.Y - c.Position.Y
	if absY < 0 {
		absY = -absY
	}
	r := absX + absY

	return r
}
