package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
)

func TestEightEntitiesFighting(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	registry.RegisterWorld(world)

	world.AddPlayer()
	world.AddPlayer()

	unitData := []byte("{\"knight\":[{\"x\":5,\"y\":3},{\"x\":4,\"y\":6},{\"x\":3,\"y\":6},{\"x\":5,\"y\":6}]}")
	unitData2 := []byte("{\"knight\":[{\"x\":10,\"y\":15},{\"x\":4,\"y\":15},{\"x\":3,\"y\":15},{\"x\":5,\"y\":15}]}")
	world.AddPlayerUnits(unitData, 0)
	world.AddPlayerUnits(unitData2, 1)

	p1 := world.ObjectPool.Components["PositionComponent"][0].(components.PositionComponent)
	p2 := world.ObjectPool.Components["PositionComponent"][1].(components.PositionComponent)
	p3 := world.ObjectPool.Components["PositionComponent"][2].(components.PositionComponent)
	p4 := world.ObjectPool.Components["PositionComponent"][3].(components.PositionComponent)
	p5 := world.ObjectPool.Components["PositionComponent"][4].(components.PositionComponent)
	p6 := world.ObjectPool.Components["PositionComponent"][5].(components.PositionComponent)
	p7 := world.ObjectPool.Components["PositionComponent"][6].(components.PositionComponent)
	p8 := world.ObjectPool.Components["PositionComponent"][7].(components.PositionComponent)

	p1.Position = engine.Vector{0, 1}
	p2.Position = engine.Vector{0, 2}
	p5.Position = engine.Vector{0, 3}
	p6.Position = engine.Vector{0, 4}
	p4.Position = engine.Vector{1, 1}
	p7.Position = engine.Vector{1, 2}
	p8.Position = engine.Vector{1, 3}
	p3.Position = engine.Vector{1, 4}

	world.ObjectPool.Components["PositionComponent"][0] = p1
	world.ObjectPool.Components["PositionComponent"][1] = p2
	world.ObjectPool.Components["PositionComponent"][2] = p3
	world.ObjectPool.Components["PositionComponent"][3] = p4
	world.ObjectPool.Components["PositionComponent"][4] = p5
	world.ObjectPool.Components["PositionComponent"][5] = p6
	world.ObjectPool.Components["PositionComponent"][6] = p7
	world.ObjectPool.Components["PositionComponent"][7] = p8

	h1 := world.ObjectPool.Components["StatsComponent"][1].(components.StatsComponent)
	h1.Health = 30

	h5 := world.ObjectPool.Components["StatsComponent"][5].(components.StatsComponent)
	h5.Health = 18

	h7 := world.ObjectPool.Components["StatsComponent"][7].(components.StatsComponent)
	h7.Health = 45

	world.ObjectPool.Components["StatsComponent"][1] = h1
	world.ObjectPool.Components["StatsComponent"][5] = h5
	world.ObjectPool.Components["StatsComponent"][7] = h7

	for world.MatchActive {
		world.Update()
	}
}

func TestTwoEntitiesFighting(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	registry.RegisterWorld(world)

	world.AddPlayer()
	world.AddPlayer()

	unitData := []byte("{\"knight\":[{\"x\":8,\"y\":9}]}")
	unitData2 := []byte("{\"knight\":[{\"x\":8,\"y\":9}]}")
	world.AddPlayerUnits(unitData, 0)
	world.AddPlayerUnits(unitData2, 1)

	h1 := world.ObjectPool.Components["StatsComponent"][0].(components.StatsComponent)
	h1.Health = 40

	h5 := world.ObjectPool.Components["StatsComponent"][1].(components.StatsComponent)
	h5.Health = 40

	world.ObjectPool.Components["StatsComponent"][0] = h1
	world.ObjectPool.Components["StatsComponent"][1] = h5

	for world.MatchActive {
		world.Update()
	}
}
