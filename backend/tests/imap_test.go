package tests

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/grid"
	"github.com/Nikola-Milovic/tog-plugin/startup"
	"testing"
)

func TestWriting(t *testing.T) {
	printImapsToFile()
}

func TestSmallerIntoBigger(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":8}]}}")

	world := CreateTestWorld(u1, u2, t)
	grid := world.Grid

	imap := grid.GetEnemyProximityImap(0)
	imap2 := startup.ProximityTemplates[2].Imap
	engine.AddIntoBiggerMap(imap2, imap,
		0, 0, 1)

	engine.PrintImapToFile(imap, "AddIntoBiggerMap", false)
}

func TestBiggerIntoSmaller(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":8}]}}")

	CreateTestWorld(u1, u2, t)
	//	grid := world.Grid

	bigger := startup.ProximityTemplates[3].Imap
	smaller := startup.ProximityTemplates[1].Imap
	engine.PrintImapToFile(smaller, "Before", false)
	engine.PrintImapToFile(bigger, "Bigger", true)
	engine.AddIntoSmallerMap(bigger, smaller, 3, 3, -1)

	engine.PrintImapToFile(smaller, "BiggerIntoSmallerMap", true)
}

func TestNormalizedMap(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":8}]}}")

	CreateTestWorld(u1, u2, t)
	//	grid := world.Grid

	bigger := startup.ProximityTemplates[3].Imap
	smaller := startup.ProximityTemplates[0].Imap
	engine.AddIntoSmallerMap(bigger, smaller, 6, 6, 2)
	engine.PrintImapToFile(smaller, "Before", false)
	smaller.Normalize()
	engine.PrintImapToFile(smaller, "Normalized", true)
}

func TestNormalizedAndInverted(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":8}]}}")

	CreateTestWorld(u1, u2, t)
	//	grid := world.Grid

	bigger := startup.ProximityTemplates[3].Imap
	smaller := startup.ProximityTemplates[0].Imap
	engine.AddIntoSmallerMap(bigger, smaller, 6, 6, 2)
	engine.PrintImapToFile(smaller, "Before", false)
	smaller.NormalizeAndInvert()
	engine.PrintImapToFile(smaller, "Inverted", true)
}

func TestSmallestValue(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":11,\"y\":1}, {\"x\":11,\"y\":2}, {\"x\":11,\"y\":1}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":14,\"y\":1}]}}")

	world := CreateTestWorld(u1, u2, t)
	g := world.Grid

	g.Update()

	workingMap := engine.NewImap(16, 16, constants.TileSize)

	posComp := world.ObjectPool.Components["PositionComponent"][0]
	pos := posComp.(components.PositionComponent).Position
	x := int(pos.X / constants.TileSize)
	y := int(pos.Y / constants.TileSize)

	engine.AddIntoSmallerMap(g.GetEnemyProximityImap(0), workingMap, x, y, -1.5)
	engine.AddIntoSmallerMap(g.GetEnemyProximityImap(1), workingMap, x, y, 2)
	engine.PrintImapToFile(workingMap, "Workingmap", false)

	x, y, _ = workingMap.GetLowestValue()
	fmt.Printf("Smallest value is at  X : %d Y : %d", x, y)
}

func TestHighestValue(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":11,\"y\":1}, {\"x\":11,\"y\":2}, {\"x\":11,\"y\":1}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":14,\"y\":1}]}}")

	world := CreateTestWorld(u1, u2, t)
	g := world.Grid

	g.Update()

	workingMap := engine.NewImap(15, 15, constants.TileSize)

	posComp := world.ObjectPool.Components["PositionComponent"][0]
	pos := posComp.(components.PositionComponent).Position
	x := int(pos.X / constants.TileSize)
	y := int(pos.Y / constants.TileSize)

	engine.AddIntoSmallerMap(g.GetEnemyProximityImap(0), workingMap, x, y, -1)
	engine.AddIntoSmallerMap(g.GetEnemyProximityImap(1), workingMap, x, y, 3)
	engine.PrintImapToFile(workingMap, "Workingmap", false)

	x, y, _ = workingMap.GetHighestCell()
	fmt.Printf("Highest value is at  X : %d Y : %d", x, y)
}

func TestTranslatingCoordsFromImapToBase(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":0,\"y\":1}, {\"x\":11,\"y\":2}, {\"x\":11,\"y\":1}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":14,\"y\":1}]}}")

	world := CreateTestWorld(u1, u2, t)
	g := world.Grid

	g.Update()

	workingMap := engine.NewImap(15, 15, constants.TileSize)

	posComp := world.ObjectPool.Components["PositionComponent"][0]
	pos := posComp.(components.PositionComponent).Position
	x := int(pos.X / constants.TileSize)
	y := int(pos.Y / constants.TileSize)
	fmt.Printf("Translated value is is at  X : %d Y : %d", x, y)

	engine.AddIntoSmallerMap(g.GetEnemyProximityImap(0), workingMap, x, y, -1)
	engine.AddIntoSmallerMap(g.GetEnemyProximityImap(1), workingMap, x, y, 3)
	engine.PrintImapToFile(workingMap, "Workingmap", false)

	lowX, lowY,_ := workingMap.GetLowestValue()
	x, y = grid.GetBaseMapCoordsFromSectionImapCoords(x, y, lowX, lowY)
	fmt.Printf("Translated value is is at  X : %d Y : %d, was X: %d Y : %d", x, y, lowX, lowY)
}

func TestSubtractingMyOwnProximity(t *testing.T) { // TODO edge of map
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":4}, {\"x\":11,\"y\":2}, {\"x\":11,\"y\":1}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":14,\"y\":1}]}}")

	world := CreateTestWorld(u1, u2, t)
	g := world.Grid

	g.Update()

	workingMap := engine.NewImap(50, 50, constants.TileSize)

	tag := world.GetEntityManager().GetEntities()[0].PlayerTag

	posComp := world.ObjectPool.Components["PositionComponent"][0]
	movComp := world.ObjectPool.Components["MovementComponent"][0].(components.MovementComponent)
	pos := posComp.(components.PositionComponent).Position
	x := int(pos.X / constants.TileSize)
	y := int(pos.Y / constants.TileSize)

	engine.AddIntoSmallerMap(g.GetProximityImaps()[tag], workingMap, x, y, 1)
	engine.PrintImapToFile(g.GetProximityImaps()[tag], "ProximityImap", false)
	engine.PrintImapToFile(grid.GetProximityTemplate(movComp.MovementSpeed).Imap, "ProximityTemplate", true)
	engine.AddIntoBiggerMap(grid.GetProximityTemplate(movComp.MovementSpeed).Imap, workingMap, 25, 25, -1)
	engine.PrintImapToFile(workingMap, "0's", true)
}