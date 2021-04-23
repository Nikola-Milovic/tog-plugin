package tests

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/grid"
	"github.com/Nikola-Milovic/tog-plugin/startup"
	"reflect"
	"testing"
)

func Test_getProximityTemplate(t *testing.T) {
	type args struct {
		speed float32
	}
	tests := []struct {
		name string
		args args
		want *engine.ImapTemplate
	}{
		{name: "0", args: args{0}, want: startup.ProximityTemplates[0]},
		{name: "Slowest", args: args{4}, want: startup.ProximityTemplates[1]},
		{name: "0", args: args{2}, want: startup.ProximityTemplates[1]},
		{name: "Slowest", args: args{5}, want: startup.ProximityTemplates[2]},
		{name: "Slowest", args: args{9}, want: startup.ProximityTemplates[3]},
		{name: "Slowest", args: args{7}, want: startup.ProximityTemplates[2]},
		{name: "0", args: args{12}, want: startup.ProximityTemplates[3]},
		{name: "0", args: args{15}, want: startup.ProximityTemplates[4]},
		{name: "0", args: args{19}, want: startup.ProximityTemplates[4]},
		{name: "0", args: args{21}, want: startup.ProximityTemplates[4]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := grid.GetProximityTemplate(tt.args.speed); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getProximityTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGridUpdate(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":1,\"y\":1}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":1,\"y\":1}]}}")

	world := CreateTestWorld(u1, u2, t)
	grid := world.Grid

	grid.Update()

	engine.PrintImapToFile(grid.GetEnemyProximityImap(0), "Updated Grid 0", false)
	engine.PrintImapToFile(grid.GetEnemyProximityImap(1), "Updated Grid 1", true)
	engine.PrintImapToFile(grid.GetOccupationalMap(), "OccupationalGrid", true)
}

func TestGridUpdateOverlappedEnemies(t *testing.T) {
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

	engine.AddIntoSmallerMap(g.GetEnemyProximityImap(0), workingMap, x, y, -1)
	engine.AddIntoSmallerMap(g.GetEnemyProximityImap(1), workingMap, x, y, 0.8)
	engine.PrintImapToFile(workingMap, "Workingmap", false)
}

func TestGetWorkingMap(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":11,\"y\":1}, {\"x\":11,\"y\":2}, {\"x\":11,\"y\":1}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":14,\"y\":1}]}}")

	world := CreateTestWorld(u1, u2, t)
	g := world.Grid

	g.Update()

	speedIn2Seconds := 2 * 7 * constants.TickRate
	workingMap := g.GetWorkingMap(speedIn2Seconds, speedIn2Seconds)

	engine.PrintImapToFile(workingMap, "Workingmap", false)
}

func TestAddingMapsToWorkingMap(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":11,\"y\":1}, {\"x\":11,\"y\":2}, {\"x\":11,\"y\":1}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":14,\"y\":1}]}}")

	world := CreateTestWorld(u1, u2, t)
	g := world.Grid

	g.Update()

	posComp := world.ObjectPool.Components["PositionComponent"][0]
	pos := posComp.(components.PositionComponent).Position
	x := int(pos.X / constants.TileSize)
	y := int(pos.Y / constants.TileSize)

	speedIn2Seconds := 2 * 7 * constants.TickRate
	workingMap := g.GetWorkingMap(speedIn2Seconds, speedIn2Seconds)

	engine.AddIntoSmallerMap(g.GetOccupationalMap(), workingMap, x, y, 3)
	engine.AddIntoSmallerMap(g.GetEnemyProximityImap(0), workingMap, x, y, 1.2)
	engine.AddIntoSmallerMap(g.GetEnemyProximityImap(1), workingMap, x, y, -1)
	workingMap.NormalizeAndInvert()
	engine.AddIntoSmallerMap(startup.InterestTemplates[1].Imap, workingMap, x, y, 1)
	x, y, _ = workingMap.GetHighestCell()

	engine.PrintImapToFile(workingMap, fmt.Sprintf("X: %d Y: %d", x, y), false)
}

func TestGetInterestTemplate(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":11,\"y\":1}, {\"x\":11,\"y\":2}, {\"x\":11,\"y\":1}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":14,\"y\":1}]}}")

	world := CreateTestWorld(u1, u2, t)
	g := world.Grid

	g.Update()

	//	posComp := world.ObjectPool.Components["PositionComponent"][0]
	//	pos := posComp.(components.PositionComponent).Position
	//x := int(pos.X / constants.TileSize)
	//y := int(pos.Y / constants.TileSize)

	speedIn2Seconds := 2 * 4 * constants.TickRate
	//workingMap := g.GetWorkingMap(speedIn2Seconds, speedIn2Seconds)

	interest := g.GetInterestTemplate(speedIn2Seconds)

	engine.PrintImapToFile(interest, fmt.Sprintf("Interest Template for speed 4"), false)
}
