package engine

import (
	"math"
)

// Grid represents the grid in the game, it is used for movement, pathfinding, abilities, attacks
// cell size is 32x32
type Grid struct {
	CellSize  int
	MaxWidth  int
	MaxHeight int
	cells     map[int]map[int]*Cell
}

// CreateGrid initializes the grid with the basic data and returns a pointer to it
func CreateGrid() *Grid { //TODO: check if should be <=
	g := Grid{}

	g.CellSize = 32
	g.MaxWidth = 800
	g.MaxHeight = 512

	g.cells = make(map[int]map[int]*Cell)

	for x := 0; x < g.MaxWidth/g.CellSize; x++ {
		for y := 0; y < g.MaxHeight/g.CellSize; y++ {
			g.setCell(&Cell{Position: Vector{X: x, Y: y}, isOccupied: false, grid: &g}, x, y)
		}
	}

	return &g
}

//Update the grid, called every tick
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

func (g *Grid) setCell(c *Cell, x, y int) {
	if g.cells[x] == nil {
		g.cells[x] = map[int]*Cell{}
	}

	g.cells[x][y] = c
}

//IsCellTaken
func (g *Grid) IsCellTaken(pos Vector) bool {
	if pos.X < 0 || pos.Y < 0 || pos.X > g.MaxWidth || pos.Y > g.MaxHeight {
		return false
	}
	return g.cells[pos.X][pos.Y].isOccupied
}

//CellAt returns a pointer to the cell at the give position and boolean indicating whether it was found or not
func (g *Grid) CellAt(pos Vector) (*Cell, bool) {
	cell, ok := g.cells[pos.X][pos.Y]

	if !ok {
		return nil, false
	}

	return cell, true
}

//GetPath returns a path from argument1 to argument2, a distance and a boolean to indicate whether or not a path was found
// if not found, returns empty slice for path
func (g *Grid) GetPath(from Vector, to Vector) (path []Vector, distance int, found bool) {
	if start, ok := g.CellAt(from); !ok {
		return []Vector{}, -1, false
	} else if end, ok := g.CellAt(to); !ok {
		return []Vector{}, -1, false
	} else {
		return Path(*start, *end)
	}
}

//GetNeighbours returns left, right, up, down neighbouring cell, ignores occupied
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

//OccupyCell indicates that cell CANNOT be occupied and is taken now by another entity
func (g *Grid) OccupyCell(coordinates Vector, id string) {
	g.cells[coordinates.X][coordinates.Y].isOccupied = true
	g.cells[coordinates.X][coordinates.Y].OccupiedID = id
}

//ReleaseCell indicates that cell can be occupied and is free now
func (g *Grid) ReleaseCell(coordinates Vector) {
	g.cells[coordinates.X][coordinates.Y].isOccupied = false
	g.cells[coordinates.X][coordinates.Y].OccupiedID = ""
	g.cells[coordinates.X][coordinates.Y].Flag.OccupiedInSteps = -1
}

//GetSurroundingTiles returns all 8 surrounding ciles, ignores occupied
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

//GetDistance returns the distance the way unit would travel, without diagonals
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

//GetDistanceIncludingDiagonal returns distance along with diagonals
func (g *Grid) GetDistanceIncludingDiagonal(c1 Vector, c2 Vector) int {

	r := math.Max(math.Abs(float64(c1.X-c2.X)), math.Abs(float64(c1.Y-c2.Y)))

	return int(r)
}
