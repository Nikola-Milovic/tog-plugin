package grid

import "C"
import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
	"github.com/Nikola-Milovic/tog-plugin/startup"
)

// Grid is
type Grid struct { // TODO maybe pointers
	MaxWidth         int
	MaxHeight        int
	UpdateInterval   int
	world            *game.World
	proximityIMaps   []*engine.Imap
	occupationalIMap *engine.Imap
	workingMap       *engine.Imap
	goalMap          *engine.Imap
}

var tileSize = float32(4)

// CreateGrid initializes
func CreateGrid(w *game.World) *Grid {
	g := Grid{}

	g.MaxWidth = constants.MapWidth
	g.MaxHeight = constants.MapHeight

	TileSize := constants.TileSize

	g.occupationalIMap = engine.NewImap(constants.MapWidth/constants.TileSize, constants.MapHeight/constants.TileSize, TileSize)

	g.goalMap = engine.NewImap(constants.MapWidth/constants.TileSize, constants.MapHeight/constants.TileSize, TileSize)

	g.workingMap = engine.NewImap((constants.MapWidth)/constants.TileSize, (constants.MapHeight)/constants.TileSize, TileSize)

	g.proximityIMaps = make([]*engine.Imap, 2)

	g.proximityIMaps[0] = engine.NewImap(constants.MapWidth/constants.TileSize, constants.MapHeight/constants.TileSize, TileSize)
	g.proximityIMaps[1] = engine.NewImap(constants.MapWidth/constants.TileSize, constants.MapHeight/constants.TileSize, TileSize)

	g.world = w
	return &g
}

func (g *Grid) GetWorkingMap(width, height int) *engine.Imap {
	if width > cap(g.workingMap.Grid) {
		width = cap(g.workingMap.Grid)
	}
	if height > cap(g.workingMap.Grid[0]) {
		height = cap(g.workingMap.Grid[0])
	}

	g.workingMap.Grid = g.workingMap.Grid[:width]

	g.workingMap.Width = width
	g.workingMap.Height = height

	for index := range g.workingMap.Grid {
		g.workingMap.Grid[index] = g.workingMap.Grid[index][:height]
	}

	g.workingMap.Clear()

	return g.workingMap
}

func (g *Grid) GetEnemyProximityImap(tag int) *engine.Imap {
	opposingTag := 0
	if tag == 0 {
		opposingTag = 1
	}

	return g.proximityIMaps[opposingTag]
}

func (g *Grid) GetOccupationalMap() *engine.Imap {
	return g.occupationalIMap
}

func (g *Grid) GetGoalMap() *engine.Imap {
	return g.goalMap
}

func (g *Grid) GetProximityImaps() []*engine.Imap {
	return g.proximityIMaps
}

func (g *Grid) GetInterestTemplate(size int) *engine.Imap {
	for _, templ := range startup.InterestTemplates {
		if templ.Radius >= size {
			return templ.Imap
		}
	}

	return startup.InterestTemplates[len(startup.InterestTemplates)-1].Imap
}

//Update the grid
func (g *Grid) Update() {
	//if g.world.Tick % g.UpdateInterval != 0 {
	//	return
	//}
	return

	for _, temp := range g.proximityIMaps {
		temp.Clear()
	}

	g.occupationalIMap.Clear()

	if g.world.Tick%5 == 0 {
		g.goalMap.Clear()
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

		proxMap := g.proximityIMaps[ent.PlayerTag]

		proxTemplate := GetProximityTemplate(movementComp.MovementSpeed)

		x, y := GlobalCordToTiled(posComp.Position)
		//if movementComp.Velocity.X > 0 {
		//	fX, fY := GlobalCordToTiled(posComp.Position.Add(movementComp.Velocity))
		//	engine.AddMaps(proxTemplate.Imap, proxMap, x, y, 0.5)
		//	engine.AddMaps(proxTemplate.Imap, proxMap, fX, fY, 0.5)
		//} else {
		engine.AddIntoBiggerMap(proxTemplate.Imap, proxMap, x, y, 1)
		//	}

		sizeTemplate := GetSizeTemplate(posComp.Radius)
		engine.AddIntoBiggerMap(sizeTemplate.Imap, g.occupationalIMap, x, y, 1.0)
	}
}

func (g *Grid) IsPositionFree(pos math.Vector, bbox math.Vector) bool {
	x, y := GlobalCordToTiled(pos)
	bboxX := int(bbox.X / tileSize)
	bboxY := int(bbox.Y / tileSize)

	xStart := math.Max(0, x-bboxX)
	xEnd := math.Min(g.MaxWidth, x+bboxX)
	yStart := math.Max(0, y-bboxY)
	yEnd := math.Min(g.MaxHeight, y+bboxY)

	om := g.occupationalIMap

	for xPos := xStart; xPos < xEnd; xPos++ {
		for yPos := yStart; yPos < yEnd; yPos++ {
			if om.Grid[xPos][yPos] > 0.0 {
				return false
			}
		}
	}

	return true
}

func GlobalCordToTiled(pos math.Vector) (x, y int) {
	return int(pos.X / constants.TileSize), int(pos.Y / constants.TileSize)
}

func GetBaseMapCoordsFromSectionImapCoords(baseCenterX, baseCenterY, x, y int) (newX, newY int) {
	adaptedX := 0
	adaptedY := 0

	adaptedX = math.Max(0, math.Min(constants.MapWidth, baseCenterX+x))
	adaptedY = math.Max(0, math.Min(constants.MapHeight, baseCenterY+y))

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

func GetSizeTemplate(radius float32) *engine.ImapTemplate {
	size := constants.StandardSize
		switch radius {
		case 16:
			size = constants.StandardSize
		case 10:
			size = constants.SmallSize
		}

	return startup.SizeTemplates[size]
}
