package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

func TestEightEntitiesFighting(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := CreateTestWorld(p1Units, p2Units, t)

	p1 := world.ObjectPool.Components["PositionComponent"][0].(components.PositionComponent)
	p2 := world.ObjectPool.Components["PositionComponent"][1].(components.PositionComponent)
	p3 := world.ObjectPool.Components["PositionComponent"][2].(components.PositionComponent)
	p4 := world.ObjectPool.Components["PositionComponent"][3].(components.PositionComponent)
	p5 := world.ObjectPool.Components["PositionComponent"][4].(components.PositionComponent)
	p6 := world.ObjectPool.Components["PositionComponent"][5].(components.PositionComponent)
	p7 := world.ObjectPool.Components["PositionComponent"][6].(components.PositionComponent)
	p8 := world.ObjectPool.Components["PositionComponent"][7].(components.PositionComponent)

	p1.Position = engine.Vector{0, 5}
	p2.Position = engine.Vector{10, 30}
	p5.Position = engine.Vector{60, 80}
	p6.Position = engine.Vector{90, 40}
	p4.Position = engine.Vector{190, 35}
	p7.Position = engine.Vector{60, 50}
	p8.Position = engine.Vector{10, 30}
	p3.Position = engine.Vector{150, 40}

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
