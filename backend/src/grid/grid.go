package grid

import "github.com/Nikola-Milovic/tog-plugin/src/constants"

type Grid struct {
	tilesize  int
	maxWidth  int
	maxHeight int
	cells     map[int]map[int]*Cell
}

func (g *Grid) initializeGrid() { //TODO: check if should be <=
	for x := 0; x < g.maxWidth/g.tilesize; x++ {
		for y := 0; y < g.maxHeight/g.tilesize; y++ {
			g.cells[x][y] = &Cell{Position: constants.V2{X: x, Y: y}, isOccupied: false}
		}
	}
}

func (g *Grid) isCellTaken(x int, y int) bool {
	return g.cells[x][y].isOccupied
}

func (g *Grid) CellAt(x int, y int) *Cell {
	if g.cells[x] == nil {
		return nil
	}
	return g.cells[x][y]
}
