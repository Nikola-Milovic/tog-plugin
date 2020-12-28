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

	world := CreateTestWorld(Lemi1Units, Lemi2Units, t)

	em := world.EntityManager

	em.Entities[2].Active = false
	idToBeRemoved := em.Entities[2].ID

	em.Systems[0].Update()

	if em.Entities[2].ID == idToBeRemoved {
		t.Errorf("Id %s should have been removed", idToBeRemoved)
	}
}

func TestDeathSystemRemoveMultiple(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := CreateTestWorld(Lemi1Units, Lemi2Units, t)

	em := world.EntityManager

	em.Entities[2].Active = false
	em.Entities[5].Active = false
	em.Entities[6].Active = false
	em.Entities[7].Active = false

	em.Systems[0].Update()

	if len(em.Entities) != 4 {
		t.Errorf("After removing 4 entities, there should be 4 entities left, but instead got %v", len(em.Entities))
	}

}
