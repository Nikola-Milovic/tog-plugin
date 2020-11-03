package grid

import "github.com/Nikola-Milovic/tog-plugin/src/constants"

type Grid struct {
	tilesize  int
	maxWidth  int
	maxHeight int
	cells     map[int]map[int]Cell
}

func (g *Grid) initializeGrid() { //TODO: check if should be <=
	for x := 0; x < g.maxWidth/g.tilesize; x++ {
		for y := 0; y < g.maxHeight/g.tilesize; y++ {
			g.cells[x][y] = Cell{Position: constants.V2{X: x, Y: y}, isOccupied: false}
		}
	}
}

func (g *Grid) IsCellTaken(x int, y int) bool {
	if x < 0 || y < 0 || x > g.maxWidth || y > g.maxHeight {
		return false
	}
	return g.cells[x][y].isOccupied
}

func (g *Grid) CellAt(x int, y int) Cell {
	return g.cells[x][y]
}

func (g *Grid) GetNeighbours(x int, y int) []Cell {
	neighbours := make([]Cell, 0, 4)

	if g.cells[x-1][y].isOccupied {
		neighbours = append(neighbours, g.cells[x-1][y])
	}
	if g.cells[x+1][y].isOccupied {
		neighbours = append(neighbours, g.cells[x-1][y])
	}
	if g.cells[x][y-1].isOccupied {
		neighbours = append(neighbours, g.cells[x-1][y])
	}
	if g.cells[x][y+1].isOccupied {
		neighbours = append(neighbours, g.cells[x-1][y])
	}

	return neighbours
}
