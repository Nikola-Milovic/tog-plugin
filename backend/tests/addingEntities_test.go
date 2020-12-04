package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

func TestSingleEntityCreation(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/test/singleUnitTest.json")
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	world.EntityManager.AddEntity(data)

	println(world.EntityManager.Entities[0].Name)

	//Test if after adding an entity the length of Entities, so we aren't wasting loops
	if len(world.EntityManager.Entities) != 1 {
		t.Errorf("Added 1 entity, expected length is 1, got %v", len(world.EntityManager.Entities))
	}
	//Check if all components are registered
	if len(world.ObjectPool.Components) != len(data["Components"].(map[string]interface{})) {
		t.Errorf("%v components should be registered,"+
			"but got %v", len(data["Components"].(map[string]interface{})), len(world.ObjectPool.Components))
	}

	//Check if component size is correct
	if len(world.ObjectPool.Components["MovementComponent"]) != 1 {
		t.Errorf("Added 1 entities, expected component length 1, got %v", len(world.ObjectPool.Components["MovementComponent"]))
	}
}

func TestCorrectMovementComponentData(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/test/singleUnitTest.json")
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	world.EntityManager.AddEntity(data)

	mComp := world.ObjectPool.Components["MovementComponent"][0].(game.MovementComponent)

	if mComp.Speed != constants.MovementSpeedFast {
		t.Errorf("Expected fast speed %v got %v", constants.MovementSpeedFast, mComp.Speed)
	}
}

func TestMultipleEntityCreation(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/test/singleUnitTest.json")
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	entityNum := 5

	for i := 0; i < entityNum; i++ {
		world.EntityManager.AddEntity(data)
	}
	//Test if after adding an entity the length of Entities, so we aren't wasting loops
	if len(world.EntityManager.Entities) != entityNum {
		t.Errorf("Added %v entities, expected length is %v, got %v", entityNum, entityNum, len(world.EntityManager.Entities))
	}

	if len(world.ObjectPool.Components["MovementComponent"]) != entityNum {
		t.Errorf("Added %v entities, expected component length %v, got %v", entityNum, entityNum, len(world.ObjectPool.Components["MovementComponent"]))
	}
}

func TestAddDifferentEntities(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/test/twoUnitsTest.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	world.EntityManager.AddEntity(data[0])
	world.EntityManager.AddEntity(data[1])

	if len(world.EntityManager.Entities) != 2 {
		t.Errorf("Added 2 entities, expected length is 2, got %v", len(world.EntityManager.Entities))
	}

	if world.EntityManager.Entities[0].Name != "Knight" {
		t.Errorf("Entity at index 0 name should be Knight, got %v", world.EntityManager.Entities[0].Name)
	}

	if world.EntityManager.Entities[1].Name != "Archer" {
		t.Errorf("Entity at index 1 name should be Archer, got %v", world.EntityManager.Entities[1].Name)
	}

	if world.EntityManager.Entities[0].Index != 0 {
		t.Errorf("Entity at index 0, index should be 0, got %v", world.EntityManager.Entities[0].Index)
	}

	if world.EntityManager.Entities[1].Index != 1 {
		t.Errorf("Entity at index 1, index should be 1, got %v", world.EntityManager.Entities[1].Index)
	}
}
