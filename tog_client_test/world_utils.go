package main

import (
	"encoding/json"
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/match"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
)

func CreateWorld(unitData []byte, unitData2 []byte) *game.World {
	world := game.CreateWorld()
	registry.RegisterWorld(world)

	world.AddPlayer("")
	world.AddPlayer("")

	data1 := match.PlayerReadyDataMessage{}
	if err := json.Unmarshal(unitData, &data1); err != nil {
		panic(fmt.Sprintf("Error unmarshaling unitData %s", err.Error()))

	}

	data2 := match.PlayerReadyDataMessage{}
	if err := json.Unmarshal(unitData2, &data2); err != nil {
		panic(fmt.Sprintf("Error unmarshaling unitData %s", err.Error()))
	}

	world.AddPlayerUnits(data1.UnitData, 0)
	world.AddPlayerUnits(data2.UnitData, 1)

	return world
}