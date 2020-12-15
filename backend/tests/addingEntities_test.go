// package tests

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"reflect"
// 	"testing"

// 	"github.com/Nikola-Milovic/tog-plugin/constants"
// 	"github.com/Nikola-Milovic/tog-plugin/game"
// )

// func TestSingleEntityCreation(t *testing.T) {
// 	jsonData, _ := ioutil.ReadFile("../resources/units.json")
// 	var data map[string]interface{}
// 	err := json.Unmarshal(jsonData, &data)
// 	if err != nil {
// 		t.Errorf("Couldn't unmarshal json: %e", err)
// 	}

// 	world := game.CreateWorld()

// 	world.EntityManager.AddEntity(data)

// 	println(world.EntityManager.Entities[0].ID)

// 	//Test if after adding an entity the length of Entities, so we aren't wasting loops
// 	if len(world.EntityManager.Entities) != 1 {
// 		t.Errorf("Added 1 entity, expected length is 1, got %v", len(world.EntityManager.Entities))
// 	}
// 	//Check if all components are registered
// 	if len(world.ObjectPool.Components) != len(data["Components"].(map[string]interface{})) {
// 		t.Errorf("%v components should be registered,"+
// 			"but got %v", len(data["Components"].(map[string]interface{})), len(world.ObjectPool.Components))
// 	}

// 	//Check if component size is correct
// 	if len(world.ObjectPool.Components["MovementComponent"]) != 1 {
// 		t.Errorf("Added 1 entities, expected component length 1, got %v", len(world.ObjectPool.Components["MovementComponent"]))
// 	}

// 	//Check if AI is added correctly
// 	if reflect.TypeOf(world.ObjectPool.AI["knight"]) != reflect.TypeOf(game.KnightAI{}) {
// 		t.Errorf("AI for knight should be type of %v, instead got %v", reflect.TypeOf(game.KnightAI{}), reflect.TypeOf(world.ObjectPool.AI["knight"]))
// 	}
// }

// func TestCorrectComponentValues(t *testing.T) {
// 	jsonData, _ := ioutil.ReadFile("../resources/test/singleUnitTest.json")
// 	var data map[string]interface{}
// 	err := json.Unmarshal(jsonData, &data)
// 	if err != nil {
// 		t.Errorf("Couldn't unmarshal json: %e", err)
// 	}

// 	world := game.CreateWorld()

// 	world.EntityManager.AddEntity(data)

// 	movementComponent := world.ObjectPool.Components["MovementComponent"][0].(game.MovementComponent)
// 	healthComponent := world.ObjectPool.Components["HealthComponent"][0].(game.HealthComponent)
// 	attackComponent := world.ObjectPool.Components["AttackComponent"][0].(game.AttackComponent)

// 	if movementComponent.MovementSpeed != constants.MovementSpeedFast {
// 		t.Errorf("Expected movement speed %v, got %v", constants.MovementSpeedFast, movementComponent.MovementSpeed)
// 	}

// 	if healthComponent.Health != 15 || healthComponent.MaxHealth != 15 {
// 		t.Errorf("Expected Health and MaxHealth to be %v, got %v", 15, healthComponent.Health)
// 	}

// 	if attackComponent.Damage != 4 {
// 		t.Errorf("Expected AttackDamage to be %v, got %v", 4, attackComponent.Damage)
// 	}

// 	if attackComponent.AttackSpeed != 8 {
// 		t.Errorf("Expected AttackSpeed to be %v, got %v", 8, attackComponent.AttackSpeed)
// 	}

// }

// func TestMultipleEntityCreation(t *testing.T) {
// 	jsonData, _ := ioutil.ReadFile("../resources/test/singleUnitTest.json")
// 	var data map[string]interface{}
// 	err := json.Unmarshal(jsonData, &data)
// 	if err != nil {
// 		t.Errorf("Couldn't unmarshal json: %e", err)
// 	}

// 	world := game.CreateWorld()

// 	entityNum := 5

// 	for i := 0; i < entityNum; i++ {
// 		world.EntityManager.AddEntity(data)
// 	}
// 	//Test if after adding an entity the length of Entities, so we aren't wasting loops
// 	if len(world.EntityManager.Entities) != entityNum {
// 		t.Errorf("Added %v entities, expected length is %v, got %v", entityNum, entityNum, len(world.EntityManager.Entities))
// 	}

// 	if len(world.ObjectPool.Components["MovementComponent"]) != entityNum {
// 		t.Errorf("Added %v entities, expected component length %v, got %v", entityNum, entityNum, len(world.ObjectPool.Components["MovementComponent"]))
// 	}
// }

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
package tests