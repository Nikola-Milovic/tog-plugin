package engine

import (
	"math"
)

type Grid struct {
	tilesize  int
	maxWidth  int
	maxHeight int
	cells     map[int]map[int]*Cell
}

func CreateGrid() *Grid { //TODO: check if should be <=
	g := Grid{}

	g.tilesize = 32
	g.maxWidth = 512
	g.maxHeight = 800

	g.cells = make(map[int]map[int]*Cell)

	for x := 0; x < g.maxWidth/g.tilesize; x++ {
		for y := 0; y < g.maxHeight/g.tilesize; y++ {
			g.SetCell(&Cell{Position: Vector{X: x, Y: y}, isOccupied: false, grid: &g}, x, y)
		}
	}

	return &g
}

func (g *Grid) Update() {
	for row := range g.cells {
		for column := range g.cells[row] {
			cell, _ := g.CellAt(Vector{row, column})
			if cell.Flag.OccupiedInSteps != -1 {
				cell.Flag.OccupiedInSteps--
			}
		}
	}
}

func (g *Grid) SetCell(c *Cell, x, y int) {
	if g.cells[x] == nil {
		g.cells[x] = map[int]*Cell{}
	}

	g.cells[x][y] = c
}

func (g *Grid) IsCellTaken(pos Vector) bool {
	if pos.X < 0 || pos.Y < 0 || pos.X > g.maxWidth || pos.Y > g.maxHeight {
		return false
	}
	return g.cells[pos.X][pos.Y].isOccupied
}

func (g *Grid) CellAt(pos Vector) (*Cell, bool) {
	cell, ok := g.cells[pos.X][pos.Y]

	if !ok {
		return &Cell{}, false
	}

	return cell, true
}

func (g *Grid) GetPath(from Vector, to Vector) (path []Vector, distance int, found bool) {
	if start, ok := g.CellAt(from); !ok {
		return []Vector{}, -1, false
	} else if end, ok := g.CellAt(to); !ok {
		return []Vector{}, -1, false
	} else {
		return Path(*start, *end)
	}
}

func (g *Grid) GetNeighbours(pos Vector) []*Cell {
	neighbours := make([]*Cell, 0, 4)

	if cell, ok := g.CellAt(Vector{X: pos.X - 1, Y: pos.Y}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	if cell, ok := g.CellAt(Vector{X: pos.X + 1, Y: pos.Y}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	if cell, ok := g.CellAt(Vector{X: pos.X, Y: pos.Y - 1}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	if cell, ok := g.CellAt(Vector{X: pos.X, Y: pos.Y + 1}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}

	return neighbours
}

func (g *Grid) OccupyCell(coordinates Vector) {
	g.cells[coordinates.X][coordinates.Y].isOccupied = true
}
func (g *Grid) ReleaseCell(coordinates Vector) {
	g.cells[coordinates.X][coordinates.Y].isOccupied = false
	g.cells[coordinates.X][coordinates.Y].Flag.OccupiedInSteps = -1
}

func (g *Grid) GetSurroundingTiles(pos Vector) []Vector {
	neighbours := make([]Vector, 0, 8)

	//left
	if cell, ok := g.CellAt(Vector{X: pos.X - 1, Y: pos.Y}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//right
	if cell, ok := g.CellAt(Vector{X: pos.X + 1, Y: pos.Y}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//down
	if cell, ok := g.CellAt(Vector{X: pos.X, Y: pos.Y - 1}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//up
	if cell, ok := g.CellAt(Vector{X: pos.X, Y: pos.Y + 1}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//top left
	if cell, ok := g.CellAt(Vector{X: pos.X - 1, Y: pos.Y - 1}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//top right
	if cell, ok := g.CellAt(Vector{X: pos.X + 1, Y: pos.Y - 1}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//bottom left
	if cell, ok := g.CellAt(Vector{X: pos.X - 1, Y: pos.Y + 1}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//bottom right
	if cell, ok := g.CellAt(Vector{X: pos.X + 1, Y: pos.Y + 1}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}

	return neighbours
}

func (g *Grid) GetDistance(c1 Vector, c2 Vector) int {
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

func (g *Grid) GetDistanceIncludingDiagonal(c1 Vector, c2 Vector) int {

	r := math.Max(math.Abs(float64(c1.X-c2.X)), math.Abs(float64(c1.Y-c2.Y)))

	return int(r)
}
