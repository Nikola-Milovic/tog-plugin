package grid

import "C"
import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

// Grid is
type Grid struct { // TODO maybe pointers
	MaxWidth       int
	MaxHeight      int
	tiles          map[int]map[int]*Tile
	world          *game.World
	proximityIMaps []engine.Imap
}

var MapWidth = 800
var MapHeight = 512
var TileSize = 4

// CreateGrid initializes
func CreateGrid(w *game.World) *Grid {
	g := Grid{}

	g.MaxWidth = MapWidth
	g.MaxHeight = MapHeight

	g.proximityIMaps = make([]engine.Imap, 2)

	g.proximityIMaps[0] = engine.NewImap(MapWidth/constants.TileSize, MapHeight/constants.TileSize, TileSize)
	g.proximityIMaps[1] = engine.NewImap(MapWidth/constants.TileSize, MapHeight/ constants.TileSize, TileSize)

	g.world = w
	return &g
}

func (g *Grid) GetEnemyProximityImap(tag int) engine.Imap {
	return g.proximityIMaps[tag]
}

//Update the grid
func (g *Grid) Update() {
}
