package grid

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

type Grid struct {
	tilesize  int
	maxWidth  int
	maxHeight int
	cells     map[int]map[int]Cell
}

func (g *Grid) InitializeGrid() { //TODO: check if should be <=

	fmt.Println("Grid intialized")

	g.tilesize = 32
	g.maxWidth = 512
	g.maxHeight = 800

	g.cells = make(map[int]map[int]Cell)

	for x := 0; x < g.maxWidth/g.tilesize; x++ {
		for y := 0; y < g.maxHeight/g.tilesize; y++ {
			g.SetCell(Cell{Position: constants.V2{X: x, Y: y}, isOccupied: false, grid: g}, x, y)
		}
	}
}

func (g *Grid) SetCell(c Cell, x, y int) {
	if g.cells[x] == nil {
		g.cells[x] = map[int]Cell{}
	}

	g.cells[x][y] = c
}

func (g *Grid) IsCellTaken(x int, y int) bool {
	if x < 0 || y < 0 || x > g.maxWidth || y > g.maxHeight {
		return false
	}
	return g.cells[x][y].isOccupied
}

func (g *Grid) CellAt(x int, y int) (Cell, bool) {
	cell, ok := g.cells[x][y]

	if !ok {
		return Cell{}, false
	}

	return cell, true
}

func (g *Grid) GetPath(from constants.V2, to constants.V2) (path []constants.V2, distance int, found bool) {
	if start, ok := g.CellAt(from.X, from.Y); !ok {
		return []constants.V2{}, -1, false
	} else if end, ok := g.CellAt(to.X, to.Y); !ok {
		return []constants.V2{}, -1, false
	} else {
		return Path(start, end)
	}
}

func (g *Grid) GetNeighbours(x int, y int) []Cell {
	neighbours := make([]Cell, 0, 4)

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

func (g *Grid) GetDistance(c1 constants.V2, c2 constants.V2) int {
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
