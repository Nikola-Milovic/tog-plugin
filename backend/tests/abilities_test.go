package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

func TestWillThrowSmite(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	unitData := []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":9,\"y\":10}]}}")
	unitData2 := []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":9,\"y\":10}]}}")
	world := CreateTestWorld(unitData, unitData2, t)

	archAtk := world.ObjectPool.Components["AttackComponent"][1].(components.AttackComponent)

	world.ObjectPool.Components["AttackComponent"][1] = archAtk

	h1 := world.ObjectPool.Components["StatsComponent"][0].(components.StatsComponent)
	h1.Health = 40

	h5 := world.ObjectPool.Components["StatsComponent"][1].(components.StatsComponent)
	h5.Health = 40

	world.ObjectPool.Components["StatsComponent"][0] = h1
	world.ObjectPool.Components["StatsComponent"][1] = h5

	for i := 0; i < 80; i++ {
		world.Update()
	}
}

func TestSummonAbility(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	var units1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"gob_beast_master\":[{\"x\":9,\"y\":10}]}}")
	var units2 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":1,\"y\":10}]}}")
	world := CreateTestWorld(units1, units2, t)

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
