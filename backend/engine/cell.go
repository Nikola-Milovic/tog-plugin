package engine

type Cell struct {
	Position   Vector
	isOccupied bool
	grid       *Grid
	Index      int
}

//instead of direct neightbours https://gamedevelopment.tutsplus.com/tutorials/how-to-speed-up-a-pathfinding-with-the-jump-point-search-algorithm--gamedev-5818

func (c Cell) PathNeighbors() []*Cell {
	return c.grid.GetNeighbours(c.Position.x, c.Position.y)
}

func (c Cell) PathEstimatedCost(co Cell) int {
	absX := co.Position.x - c.Position.x
	if absX < 0 {
		absX = -absX
	}
	absY := co.Position.y - c.Position.y
	if absY < 0 {
		absY = -absY
	}
	r := absX + absY

	return r
}
