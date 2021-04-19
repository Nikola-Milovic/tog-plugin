package grid

import "C"
import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/startup"
)

// Grid is
type Grid struct { // TODO maybe pointers
	MaxWidth          int
	MaxHeight         int
	UpdateInterval    int
	tiles             map[int]map[int]*Tile
	world             *game.World
	ProximityIMaps    []*engine.Imap
	OccupationalIMaps []*engine.Imap
}

var MapWidth = 800
var MapHeight = 512
var TileSize = 4

// CreateGrid initializes
func CreateGrid(w *game.World) *Grid {
	g := Grid{}

	g.MaxWidth = MapWidth
	g.MaxHeight = MapHeight

	g.ProximityIMaps = make([]*engine.Imap, 2)
	g.OccupationalIMaps = make([]*engine.Imap, 2)

	g.ProximityIMaps[0] = engine.NewImap(MapWidth/constants.TileSize, MapHeight/constants.TileSize, TileSize)
	g.ProximityIMaps[1] = engine.NewImap(MapWidth/constants.TileSize, MapHeight/constants.TileSize, TileSize)

	g.OccupationalIMaps[0] = engine.NewImap(MapWidth/constants.TileSize, MapHeight/constants.TileSize, TileSize)
	g.OccupationalIMaps[1] = engine.NewImap(MapWidth/constants.TileSize, MapHeight/constants.TileSize, TileSize)

	g.world = w
	return &g
}

func (g *Grid) GetEnemyProximityImap(tag int) *engine.Imap {
	opposingTag := 0
	if tag == 0 {
		opposingTag = 1
	}

	return g.ProximityIMaps[opposingTag]
}

func (g *Grid) GetEnemyOccupationalImap(tag int) *engine.Imap {
	opposingTag := 0
	if tag == 0 {
		opposingTag = 1
	}

	return g.OccupationalIMaps[opposingTag]
}

//Update the grid
func (g *Grid) Update() {
	//if g.world.Tick % g.UpdateInterval != 0 {
	//	return
	//}

	for _, temp := range g.ProximityIMaps {
		temp.Clear()
	}

	entities := g.world.EntityManager.GetEntities()
	posComps := g.world.ObjectPool.Components["PositionComponent"]
	movementComps := g.world.ObjectPool.Components["MovementComponent"]

	for _, ent := range entities {
		if !ent.Active {
			continue
		}

		posComp := posComps[ent.Index].(components.PositionComponent)
		movementComp := movementComps[ent.Index].(components.MovementComponent)

		proxMap := g.ProximityIMaps[ent.PlayerTag]

		proxTemplate := GetProximityTemplate(movementComp.MovementSpeed)

		x, y := globalCordToTiled(posComp.Position)
		engine.AddIntoBiggerMap(proxTemplate.Imap, proxMap, x, y, 1)

		sizeTemplate := GetSizeTemplate(posComp.BoundingBox)
		engine.AddIntoBiggerMap(sizeTemplate.Imap, g.OccupationalIMaps[ent.PlayerTag], x,y, 1)
	}
}

func globalCordToTiled(pos engine.Vector) (x, y int) {
	return int(pos.X / constants.TileSize), int(pos.Y / constants.TileSize)
}

func GetBaseMapCoordsFromSectionImapCoords(baseCenterX, baseCenterY, x, y int) (newX,newY int){
	adaptedX := 0
	adaptedY := 0

	adaptedX = engine.Max(0, engine.Min(MapWidth, baseCenterX + x))
	adaptedY = engine.Max(0, engine.Min(MapHeight, baseCenterY + y))

	return adaptedX, adaptedY
}

func GetProximityTemplate(speed float32) *engine.ImapTemplate {
	spd := int(speed)
	for i, templ := range startup.ProximityTemplates {
		if spd <= templ.Radius {
			return startup.ProximityTemplates[i]
		}
	}
	return startup.ProximityTemplates[len(startup.ProximityTemplates)-1]
}

func GetSizeTemplate(bbox engine.Vector) *engine.ImapTemplate {

	size := constants.StandardSize

	if bbox.X == bbox.Y {
		switch bbox.X {
		case 32:
			size = constants.StandardSize
		case 20:
			size = constants.SmallSize
		}
	}

	return startup.SizeTemplates[size]
}
