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

	for y := 0; y < 13; y++ {
		cell := g.cells[7][y]
		cell.isOccupied = true
		g.cells[7][y] = cell
	}

	return &g
}

func (g *Grid) SetCell(c *Cell, x, y int) {
	if g.cells[x] == nil {
		g.cells[x] = map[int]*Cell{}
	}

	g.cells[x][y] = c
}

func (g *Grid) IsCellTaken(x int, y int) bool {
	if x < 0 || y < 0 || x > g.maxWidth || y > g.maxHeight {
		return false
	}
	return g.cells[x][y].isOccupied
}

func (g *Grid) CellAt(x int, y int) (*Cell, bool) {
	cell, ok := g.cells[x][y]

	if !ok {
		return &Cell{}, false
	}

	return cell, true
}

func (g *Grid) GetPath(from Vector, to Vector) (path []Vector, distance int, found bool) {
	if start, ok := g.CellAt(from.X, from.Y); !ok {
		return []Vector{}, -1, false
	} else if end, ok := g.CellAt(to.X, to.Y); !ok {
		return []Vector{}, -1, false
	} else {
		return Path(*start, *end)
	}
}

func (g *Grid) GetNeighbours(x int, y int) []*Cell {
	neighbours := make([]*Cell, 0, 4)

	if cell, ok := g.CellAt(x-1, y); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	if cell, ok := g.CellAt(x+1, y); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	if cell, ok := g.CellAt(x, y-1); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	if cell, ok := g.CellAt(x, y+1); ok {
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
}

func (g *Grid) GetSurroundingTiles(x int, y int) []Vector {
	neighbours := make([]Vector, 0, 8)

	//left
	if cell, ok := g.CellAt(x-1, y); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//right
	if cell, ok := g.CellAt(x+1, y); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//down
	if cell, ok := g.CellAt(x, y-1); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//up
	if cell, ok := g.CellAt(x, y+1); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//top left
	if cell, ok := g.CellAt(x-1, y-1); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//top right
	if cell, ok := g.CellAt(x+1, y-1); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//bottom left
	if cell, ok := g.CellAt(x-1, y+1); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//bottom right
	if cell, ok := g.CellAt(x+1, y+1); ok {
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