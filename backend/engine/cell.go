package engine

type Cell struct {
	Position   Vector
	isOccupied bool
	Flag       MovementFlag
	grid       *Grid
	Index      int
}

type MovementFlag struct {
	OccupiedInSteps int
}

//instead of direct neightbours https://gamedevelopment.tutsplus.com/tutorials/how-to-speed-up-a-pathfinding-with-the-jump-point-search-algorithm--gamedev-5818

func (c Cell) PathNeighbors() []*Cell {
	return c.grid.GetNeighbours(c.Position)
}

func (c *Cell) FlagCell(inSteps int) {
	c.Flag.OccupiedInSteps = inSteps
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
