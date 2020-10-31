package grid

import "github.com/Nikola-Milovic/tog-plugin/src/constants"

type Grid struct {
	maxWidth  int
	maxHeight int
	cells     [][]Cell
}

func (g *Grid) isCellTaken(cellToCheck constants.V2) int {
	if cellToCheck.X > g.maxWidth || cellToCheck.Y > g.maxHeight || cellToCheck.X < 0 || cellToCheck.Y < 0 {
		return -1
	}
}
