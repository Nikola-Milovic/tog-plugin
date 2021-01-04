package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/match"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
)

var Lemi1Units = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":9,\"y\":10},{\"x\":9,\"y\":9},{\"x\":9,\"y\":8},{\"x\":9,\"y\":7}]}}")
var Lemi2Units = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":9,\"y\":10},{\"x\":9,\"y\":9},{\"x\":9,\"y\":8},{\"x\":9,\"y\":7}]}}")

func CreateTestWorld(unitData []byte, unitData2 []byte, testing *testing.T) *game.World {
	world := game.CreateWorld()
	registry.RegisterWorld(world)

	world.AddPlayer("")
	world.AddPlayer("")

	data1 := match.PlayerReadyDataMessage{}
	if err := json.Unmarshal(unitData, &data1); err != nil {
		fmt.Printf("Error unmarshaling unitData %s", err.Error())
		testing.FailNow()
	}

	data2 := match.PlayerReadyDataMessage{}
	if err := json.Unmarshal(unitData2, &data2); err != nil {
		fmt.Printf("Error unmarshaling unitData %s", err.Error())
		testing.FailNow()
	}

	world.AddPlayerUnits(data1.UnitData, 0)
	world.AddPlayerUnits(data2.UnitData, 1)

	return world
}
