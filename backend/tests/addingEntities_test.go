package tests

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	ai "github.com/Nikola-Milovic/tog-plugin/game/AI"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/match"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
)

// {"knight":[{"x":5,"y":3},{"x":4,"y":6},{"x":3,"y":6}]}

func TestSingleEntityCreation(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	registry.RegisterWorld(world)

	world.AddPlayer("")

	// knight : {[Positions]}
	unitData := []byte("{\"name\":\"Lemi1\",\"units\":{\"knight\":[{\"x\":9,\"y\":10}]}}")
	data1 := match.PlayerReadyDataMessage{}
	if err := json.Unmarshal([]byte(unitData), &data1); err != nil {
	}

	world.AddPlayerUnits(data1.UnitData, 0)

	//Test if after adding an entity the length of Entities, so we aren't wasting loops
	if len(world.EntityManager.Entities) != 1 {
		t.Errorf("Added 1 entity, expected length is 1, got %v", len(world.EntityManager.Entities))
	}

	components := data[0]["Components"]

	//Check if all components are registered
	if len(world.ObjectPool.Components) != len(components.(map[string]interface{})) {
		t.Errorf("%v components should be registered,"+
			"but got %v", len(components.(map[string]interface{})), len(world.ObjectPool.Components))
	}

	//Check if component size is correct
	if len(world.ObjectPool.Components["MovementComponent"]) != 1 {
		t.Errorf("Added 1 entities, expected component length 1, got %v", len(world.ObjectPool.Components["MovementComponent"]))
	}

	//Check if AI is added correctly
	if reflect.TypeOf(world.ObjectPool.AI["knight"]) != reflect.TypeOf(ai.KnightAI{}) {
		t.Errorf("AI for knight should be type of %v, instead got %v", reflect.TypeOf(ai.KnightAI{}), reflect.TypeOf(world.ObjectPool.AI["knight"]))
	}
}

func TestCorrectComponentValues(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	registry.RegisterWorld(world)

	world.AddPlayer("")

	// knight : {[Positions]}
	unitData := []byte("{\"name\":\"Lemi1\",\"units\":{\"knight\":[{\"x\":9,\"y\":10}]}}")

	data1 := match.PlayerReadyDataMessage{}
	if err := json.Unmarshal([]byte(unitData), &data1); err != nil {
	}
	world.AddPlayerUnits(data1.UnitData, 0)

	movementComponent := world.ObjectPool.Components["MovementComponent"][0].(components.MovementComponent)
	healthComponent := world.ObjectPool.Components["StatsComponent"][0].(components.StatsComponent)
	attackComponent := world.ObjectPool.Components["AttackComponent"][0].(components.AttackComponent)

	if movementComponent.MovementSpeed != constants.MovementSpeedFast {
		t.Errorf("Expected movement speed %v, got %v", constants.MovementSpeedFast, movementComponent.MovementSpeed)
	}

	if healthComponent.Health != 30 || healthComponent.MaxHealth != 30 {
		t.Errorf("Expected Health and MaxHealth to be %v, got %v", 30, healthComponent.Health)
	}

	if attackComponent.Damage != 4 {
		t.Errorf("Expected AttackDamage to be %v, got %v", 4, attackComponent.Damage)
	}

	if attackComponent.AttackSpeed != 8 {
		t.Errorf("Expected AttackSpeed to be %v, got %v", 8, attackComponent.AttackSpeed)
	}

}

func TestMultipleEntityCreation(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}
	entityNum := 8

	world := CreateTestWorld(Lemi1Units, Lemi2Units, t)

	//Test if after adding an entity the length of Entities, so we aren't wasting loops
	if len(world.EntityManager.Entities) != entityNum {
		t.Errorf("Added %v entities, expected length is %v, got %v", entityNum, entityNum, len(world.EntityManager.Entities))
	}

	if len(world.ObjectPool.Components["MovementComponent"]) != entityNum {
		t.Errorf("Added %v entities, expected component length %v, got %v", entityNum, entityNum, len(world.ObjectPool.Components["MovementComponent"]))
	}

	if world.Players[0].NumberOfUnits != 4 || world.Players[1].NumberOfUnits != 4 {
		t.Errorf("Wrong number of units!")
	}
}
