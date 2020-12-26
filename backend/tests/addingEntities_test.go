package tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	ai "github.com/Nikola-Milovic/tog-plugin/game/AI"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
	"github.com/Nikola-Milovic/tog-plugin/startup"
)

//TestMain is here to do the setup needed before all of the tests,
//populates the UnitDataMap for tests
func TestMain(m *testing.M) {
	startup.StartUp(true)
	code := m.Run()
	os.Exit(code)
}

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

	world.AddPlayer()

	// knight : {[Positions]}
	unitData := []byte("{\"knight\":[{\"x\":1,\"y\":1}]}")
	world.AddPlayerUnits(unitData, 0)

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

	world.AddPlayer()

	// knight : {[Positions]}
	unitData := []byte("{\"knight\":[{\"x\":1,\"y\":1}]}")
	world.AddPlayerUnits(unitData, 0)

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

	world := game.CreateWorld()
	registry.RegisterWorld(world)
	world.AddPlayer()
	world.AddPlayer()

	entityNum := 5

	unitData := []byte("{\"knight\":[{\"x\":5,\"y\":3},{\"x\":4,\"y\":6},{\"x\":3,\"y\":6}]}")
	unitData2 := []byte("{\"knight\":[{\"x\":5,\"y\":1},{\"x\":5,\"y\":2}]}")
	world.AddPlayerUnits(unitData, 0)
	world.AddPlayerUnits(unitData2, 1)

	//Test if after adding an entity the length of Entities, so we aren't wasting loops
	if len(world.EntityManager.Entities) != entityNum {
		t.Errorf("Added %v entities, expected length is %v, got %v", entityNum, entityNum, len(world.EntityManager.Entities))
	}

	if len(world.ObjectPool.Components["MovementComponent"]) != entityNum {
		t.Errorf("Added %v entities, expected component length %v, got %v", entityNum, entityNum, len(world.ObjectPool.Components["MovementComponent"]))
	}

	if world.Players[0].NumberOfUnits != 3 || world.Players[1].NumberOfUnits != 2 {
		t.Errorf("Wrong number of units!")
	}
}

// func TestAddDifferentEntities(t *testing.T) {
// 	jsonData, _ := ioutil.ReadFile("../resources/test/twoUnitsTest.json")
// 	var data []map[string]interface{}
// 	err := json.Unmarshal(jsonData, &data)
// 	if err != nil {
// 		t.Errorf("Couldn't unmarshal json: %e", err)
// 	}

// 	world := game.CreateWorld()

// 	world.EntityManager.AddEntity(data[0])
// 	world.EntityManager.AddEntity(data[1])

// 	if len(world.EntityManager.Entities) != 2 {
// 		t.Errorf("Added 2 entities, expected length is 2, got %v", len(world.EntityManager.Entities))
// 	}

// 	if world.EntityManager.Entities[0].ID != "knight" {
// 		t.Errorf("Entity at index 0 name should be Knight, got %v", world.EntityManager.Entities[0].ID)
// 	}

// 	if world.EntityManager.Entities[1].ID != "archer" {
// 		t.Errorf("Entity at index 1 name should be Archer, got %v", world.EntityManager.Entities[1].ID)
// 	}

// 	if world.EntityManager.Entities[0].Index != 0 {
// 		t.Errorf("Entity at index 0, index should be 0, got %v", world.EntityManager.Entities[0].Index)
// 	}

// 	if world.EntityManager.Entities[1].Index != 1 {
// 		t.Errorf("Entity at index 1, index should be 1, got %v", world.EntityManager.Entities[1].Index)
// 	}
// }
