package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

func TestApplyDotEffectPoison(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	world.AddPlayer()
	world.AddPlayer()

	unitData := []byte("{\"knight\":[{\"x\":15,\"y\":0}]}")
	unitData2 := []byte("{\"archer\":[{\"x\":14,\"y\":0}]}")
	world.AddPlayerUnits(unitData, 0)
	world.AddPlayerUnits(unitData2, 1)

	archAtk := world.ObjectPool.Components["AttackComponent"][1].(components.AttackComponent)
	archAtk.TimeSinceLastAttack = -1000
	archAtk.AttackSpeed = 100

	world.ObjectPool.Components["AttackComponent"][1] = archAtk

	h1 := world.ObjectPool.Components["StatsComponent"][0].(components.StatsComponent)
	h1.Health = 40

	h5 := world.ObjectPool.Components["StatsComponent"][1].(components.StatsComponent)
	h5.Health = 40

	world.ObjectPool.Components["StatsComponent"][0] = h1
	world.ObjectPool.Components["StatsComponent"][1] = h5

	if archAtk.OnHit != "eff_poison" {
		t.Error("OnHit != eff_poison")
	}

	for i := 0; i < 80; i++ {
		world.Update()
	}
}
