package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/match"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
)

func TestWillThrowSmite(t *testing.T) {
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

func TestOnHitAbility(t *testing.T) {
	var units1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[{\"x\":9,\"y\":10}],\"gob_beast_master\":[]}}")
	var units2 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":1,\"y\":10}]}}")

	world := game.CreateWorld()
	registry.RegisterWorld(world)

	comps := world.UnitDataMap["archer"]["Components"].(map[string]interface{})

	comps["AbilitiesComponent"] = make([]interface{}, 1, 1)
	poisonOnhit := make(map[string]interface{})
	poisonOnhit["AbilityID"] = "ab_poison_touch"
	comps["AbilitiesComponent"].([]interface{})[0] = poisonOnhit

	world.AddPlayer("")
	world.AddPlayer("")

	data1 := match.PlayerReadyDataMessage{}
	if err := json.Unmarshal(units1, &data1); err != nil {
		fmt.Printf("Error unmarshaling unitData %s", err.Error())
		t.FailNow()
	}

	data2 := match.PlayerReadyDataMessage{}
	if err := json.Unmarshal(units2, &data2); err != nil {
		fmt.Printf("Error unmarshaling unitData %s", err.Error())
		t.FailNow()
	}

	world.AddPlayerUnits(data1.UnitData, 0)
	world.AddPlayerUnits(data2.UnitData, 1)

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
