package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestDeathSystemRemoveEntity(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := CreateTestWorld(p1Units, p2Units, t)

	em := world.EntityManager

	em.GetEntities()[2].Active = false
	idToBeRemoved := em.GetEntities()[2].ID

	em.GetSystems()[0].Update()

	if em.GetEntities()[2].ID == idToBeRemoved {
		t.Errorf("Id %d should have been removed", idToBeRemoved)
	}
}

func TestDeathSystemRemoveMultiple(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	p1Units := []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[{\"x\":6,\"y\":10}],\"knight\":[{\"x\":1,\"y\":7}, {\"x\":3,\"y\":7}, {\"x\":9,\"y\":7}, {\"x\":9,\"y\":3}]}}")
	p2Units := []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[{\"x\":9,\"y\":4}, {\"x\":9,\"y\":7}],\"knight\":[{\"x\":9,\"y\":10}]}}")

	world := CreateTestWorld(p1Units, p2Units, t)

	em := world.EntityManager

	entities := em.GetEntities()

	entities[2].Active = false
	entities[5].Active = false
	entities[6].Active = false
	entities[7].Active = false

	world.Update()
	world.Update()
	world.Update()

	if len(world.EntityManager.GetEntities()) != 4 {
		t.Errorf("After removing 4 entities, there should be 4 entities left, but instead got %v", len(entities))
	}

}
