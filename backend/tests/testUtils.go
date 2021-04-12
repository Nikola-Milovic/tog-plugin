package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/match"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
)


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
