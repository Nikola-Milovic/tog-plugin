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

var CellSize = 4

// CreateGrid initializes the grid with the basic data and returns a pointer to it
func CreateGrid() *Grid { //TODO: check if should be <=
	g := Grid{}

	g.CellSize = CellSize
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
	// for row := range g.cells {
	// 	for column := range g.cells[row] {
	// 		cell, _ := g.CellAt(Vector{row, column})
	// 		if cell.Flag.OccupiedInSteps != -1 {
	// 			cell.Flag.OccupiedInSteps--
	// 		}
	// 	}
	// }
}

func (g *Grid) IsPositionAvailable(pos Vector, boundingBox Vector) bool {
	xStart := pos.X - (boundingBox.X/(CellSize))/2
	yStart := pos.Y - (boundingBox.Y/(CellSize+2))/2
	xEnd := pos.X - (boundingBox.X/(CellSize))/2
	yEnd := pos.Y - (boundingBox.Y/(CellSize+2))/2

	for x := xStart; x < xEnd; x++ {
		for y := yStart; y < yEnd; y++ {
			if g.cells[x][y].isOccupied {
				return false
			}
		}
	}
	return true
}

// func getCellPositionFromGlobalPosition(pos Vector) Vector {
// 	cellPos := Vector{}

// 	cellPos.X = pos.X/CellSize + Min(1, pos.X%CellSize)
// 	cellPos.Y = pos.Y/CellSize + Min(1, pos.Y%CellSize)

// 	return cellPos
// }

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
func (g *Grid) GetPath(from Vector, to Vector, boundingBox Vector) (path []Vector, distance int, found bool) {
	if start, ok := g.CellAt(from); !ok {
		return []Vector{}, -1, false
	} else if end, ok := g.CellAt(to); !ok {
		return []Vector{}, -1, false
	} else {
		return Path(*start, *end, boundingBox)
	}
}

//OccupyCell marks cells that are inside the bounding box of entitity
func (g *Grid) OccupyCells(pos Vector, id string, boundingBox Vector) {
	xStart := pos.X - (boundingBox.X/(CellSize))/2
	yStart := pos.Y - (boundingBox.Y/(CellSize+2))/2
	xEnd := pos.X - (boundingBox.X/(CellSize))/2
	yEnd := pos.Y - (boundingBox.Y/(CellSize+2))/2

	for x := xStart; x < xEnd; x++ {
		for y := yStart; y < yEnd; y++ {
			g.cells[x][y].isOccupied = true
			g.cells[x][y].OccupiedID = id
		}
	}
}

//ReleaseCell cells inside this bounding box are free now
func (g *Grid) ReleaseCells(pos Vector, boundingBox Vector) {
	xStart := pos.X - (boundingBox.X/(CellSize))/2
	yStart := pos.Y - (boundingBox.Y/(CellSize+2))/2
	xEnd := pos.X - (boundingBox.X/(CellSize))/2
	yEnd := pos.Y - (boundingBox.Y/(CellSize+2))/2

	for x := xStart; x < xEnd; x++ {
		for y := yStart; y < yEnd; y++ {
			g.cells[x][y].isOccupied = false
			g.cells[x][y].OccupiedID = ""
		}
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

//GetNeighbours returns left, right, up, down neighbouring cell, ignores occupied
func (g *Grid) GetNeighboursWithBoundingBox(pos Vector, boundingBox Vector) []*Cell {
	neighbours := make([]*Cell, 0, 8)

	offset := boundingBox.X

	//left
	if cell, ok := g.CellAt(Vector{X: pos.X - offset, Y: pos.Y}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	//right
	if cell, ok := g.CellAt(Vector{X: pos.X + offset, Y: pos.Y}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	//down
	if cell, ok := g.CellAt(Vector{X: pos.X, Y: pos.Y - offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	//up
	if cell, ok := g.CellAt(Vector{X: pos.X, Y: pos.Y + offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	//top left
	if cell, ok := g.CellAt(Vector{X: pos.X - offset, Y: pos.Y - offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	//top right
	if cell, ok := g.CellAt(Vector{X: pos.X + offset, Y: pos.Y - offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	//bottom left
	if cell, ok := g.CellAt(Vector{X: pos.X - offset, Y: pos.Y + offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}
	//bottom right
	if cell, ok := g.CellAt(Vector{X: pos.X + offset, Y: pos.Y + offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell)
		}
	}

	return neighbours
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

func (g *Grid) GetSurroundingTilesWithOffset(pos Vector, offset int) []Vector {
	neighbours := make([]Vector, 0, 8)

	//left
	if cell, ok := g.CellAt(Vector{X: pos.X - offset, Y: pos.Y}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//right
	if cell, ok := g.CellAt(Vector{X: pos.X + offset, Y: pos.Y}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//down
	if cell, ok := g.CellAt(Vector{X: pos.X, Y: pos.Y - offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//up
	if cell, ok := g.CellAt(Vector{X: pos.X, Y: pos.Y + offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//top left
	if cell, ok := g.CellAt(Vector{X: pos.X - offset, Y: pos.Y - offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//top right
	if cell, ok := g.CellAt(Vector{X: pos.X + offset, Y: pos.Y - offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//bottom left
	if cell, ok := g.CellAt(Vector{X: pos.X - offset, Y: pos.Y + offset}); ok {
		if !cell.isOccupied {
			neighbours = append(neighbours, cell.Position)
		}
	}
	//bottom right
	if cell, ok := g.CellAt(Vector{X: pos.X + offset, Y: pos.Y + offset}); ok {
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
