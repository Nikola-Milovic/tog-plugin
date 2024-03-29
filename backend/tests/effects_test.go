package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

func TestApplyDotEffectPoisonAndTick(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	var units1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[{\"x\":6,\"y\":10}],\"knight\":[]}}")
	var units2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":9,\"y\":10}]}}")

	world := CreateTestWorld(units1, units2, t)

	archAtk := world.ObjectPool.Components["AttackComponent"][0].(components.AttackComponent)
	archAtk.TimeSinceLastAttack = -1000
	archAtk.AttackSpeed = 100
	archAtk.OnHit = "eff_poison"

	world.ObjectPool.Components["AttackComponent"][0] = archAtk

	h1 := world.ObjectPool.Components["StatsComponent"][0].(components.StatsComponent)
	h1.Health = 40

	h5 := world.ObjectPool.Components["StatsComponent"][1].(components.StatsComponent)
	h5.Health = 40

	world.ObjectPool.Components["StatsComponent"][0] = h1
	world.ObjectPool.Components["StatsComponent"][1] = h5

	if archAtk.OnHit != "eff_poison" {
		t.Error("OnHit != eff_poison")
	}

	for world.MatchActive {
		world.Update()
	}
}
