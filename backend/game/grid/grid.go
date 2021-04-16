package grid

import "C"
import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

// Grid represents the grid in the game, it is used for movement, pathfinding, abilities, attacks
// cell size is 32x32
type Grid struct { // TODO maybe pointers
	MaxWidth        int
	MaxHeight       int
	flowTiles       []map[int]map[int]*FlowTile
	tiles           map[int]map[int]Tile
	world           *game.World
	entityPositions [][]engine.Vector
}

var FlowTileSize = 32
var TileSize = 8
var MapWidth = 800
var MapHeight = 512

// CreateGrid initializes the grid with the basic data and returns a pointer to it
func CreateGrid(w *game.World) *Grid { //TODO: check if should be <=
	g := Grid{}

	g.MaxWidth = MapWidth
	g.MaxHeight = MapHeight

	g.world = w

	g.entityPositions = make([][]engine.Vector, 2)
	g.entityPositions[0] = make([]engine.Vector, 50)
	g.entityPositions[1] = make([]engine.Vector, 50)

	g.flowTiles = make([]map[int]map[int]*FlowTile, 2)
	g.tiles = make(map[int]map[int]Tile, MapWidth/TileSize)

	g.flowTiles[0] = make(map[int]map[int]*FlowTile, MapWidth/FlowTileSize)
	g.flowTiles[1] = make(map[int]map[int]*FlowTile, MapWidth/FlowTileSize)

	for x := 0; x <= g.MaxWidth/FlowTileSize; x++ {
		for y := 0; y <= g.MaxHeight/FlowTileSize; y++ {
			g.setFlowTile(x, y)
		}
	}

	for x := 0; x <= g.MaxWidth/TileSize; x++ {
		for y := 0; y <= g.MaxHeight/TileSize; y++ {
			g.setTile(x, y)
		}
	}

	return &g
}

//Update the grid
func (g *Grid) Update() {
	g.entityPositions[0] = g.entityPositions[0][:0]
	g.entityPositions[1] = g.entityPositions[1][:0]
	entities := g.world.GetEntityManager().GetEntities()
	posComps := g.world.GetObjectPool().Components["PositionComponent"]

	for _, ent := range entities {
		g.entityPositions[ent.PlayerTag] = append(g.entityPositions[ent.PlayerTag], posComps[ent.Index].(components.PositionComponent).Position)
	}

	for x := 0; x <= g.MaxWidth/FlowTileSize; x++ {
		for y := 0; y <= g.MaxHeight/FlowTileSize; y++ {

			curPosVector := engine.Vector{X: float32(x * 32), Y: float32(y * 32)}
			// find closest tile to this one for both players
			minDist := float32(100000.0)
			minIndx := -1
			for indx, pos := range g.entityPositions[0] {
				d := engine.GetDistanceIncludingDiagonal(pos, curPosVector)
				if d < minDist {
					minIndx = indx
					minDist = d
				}
			}

			//	fmt.Printf("CurPos : %v, enemyPos : %v,  direction %v \n", curPosVector, g.entityPositions[0][minIndx], g.entityPositions[0][minIndx].Subtract(curPosVector).Normalize())
			g.flowTiles[1][x][y].Direction = g.entityPositions[0][minIndx].Subtract(curPosVector).Normalize()

			minDist1 := float32(100000.0)
			minIndx1 := -1
			for indx, pos := range g.entityPositions[1] {
				d := engine.GetDistanceIncludingDiagonal(pos, curPosVector)
				if d < minDist1 {
					minIndx1 = indx
					minDist1 = d
				}
			}

			g.flowTiles[0][x][y].Direction = g.entityPositions[1][minIndx1].Subtract(curPosVector).Normalize()
		}
	}
}

func (g *Grid) GetDesiredDirectionAt(pos engine.Vector, tag int) engine.Vector {
	x := int(pos.X) / FlowTileSize
	y := int(pos.Y) / FlowTileSize

	return g.flowTiles[tag][x][y].Direction
}

func (g *Grid) setFlowTile(x, y int) {
	if g.flowTiles[0][x] == nil {
		g.flowTiles[0][x] = map[int]*FlowTile{}
	}
	if g.flowTiles[1][x] == nil {
		g.flowTiles[1][x] = map[int]*FlowTile{}
	}

	g.flowTiles[0][x][y] = &FlowTile{}
	g.flowTiles[1][x][y] = &FlowTile{}
}

func (g *Grid) setTile(x int, y int) {
	if g.tiles[x] == nil {
		g.tiles[x] = map[int]Tile{}
	}

	g.tiles[x][y] = Tile{}
}

//func (g *Grid) IsPositionFree(index int, positionToCheck engine.Vector, boundingBox engine.Vector) int {
//	offset := int(boundingBox.X / 2)
//
//	posX := int(positionToCheck.X) / TileSize
//	posY := int(positionToCheck.X) / TileSize
//
//	xStart := engine.Max(posX-offset, 0)
//	yStart := engine.Max(0, posY-offset)
//	xEnd := engine.Min(g.MaxWidth, posX+offset)
//	yEnd := engine.Min(g.MaxHeight, posY+offset)
//
//	for x := xStart; x < xEnd; x++ {
//		for y := yStart; y < yEnd; y++ {
//			if g.tiles[x][y].isOccupied {
//				return g.tiles[x][y].occupiedIndex
//			}
//		}
//	}
//
//	return -1
//}

func (g *Grid) IsPositionFree(index int, positionToCheck engine.Vector, boundingBox engine.Vector) int {

	x := boundingBox.X / 2
	y := boundingBox.Y / 2
	leftX := positionToCheck.X - x
	rightX := positionToCheck.X + x
	topY := positionToCheck.Y - y
	bottomY := positionToCheck.Y + y

	for idx, posComp := range g.world.ObjectPool.Components["PositionComponent"] {
		pos := posComp.(components.PositionComponent).Position
		bbox := posComp.(components.PositionComponent).BoundingBox

		if idx == index {
			continue
		}

		x1 := bbox.X / 2
		y1 := bbox.Y / 2
		leftX1 := pos.X - x1
		rightX1 := pos.X + x1
		topY1 := pos.Y - y1
		bottomY1 := pos.Y + y1

		if leftX < rightX1 && rightX > leftX1 && topY > bottomY1 && topY1 < bottomY {
			return idx
		}
	}

	return -1
}
