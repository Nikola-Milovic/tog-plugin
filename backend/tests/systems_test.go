package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
)

func TestDeathSystemRemoveEntity(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	registry.RegisterWorld(world)

	em := world.EntityManager

	world.AddPlayer()
	world.AddPlayer()

	unitData := []byte("{\"knight\":[{\"x\":5,\"y\":3},{\"x\":4,\"y\":6},{\"x\":3,\"y\":6},{\"x\":5,\"y\":6}]}")
	unitData2 := []byte("{\"knight\":[{\"x\":10,\"y\":15},{\"x\":4,\"y\":15},{\"x\":3,\"y\":15},{\"x\":5,\"y\":15}]}")
	world.AddPlayerUnits(unitData, 0)
	world.AddPlayerUnits(unitData2, 1)

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

	world := game.CreateWorld()

	registry.RegisterWorld(world)

	em := world.EntityManager

	world.AddPlayer()
	world.AddPlayer()

	unitData := []byte("{\"knight\":[{\"x\":5,\"y\":3},{\"x\":4,\"y\":6},{\"x\":3,\"y\":6},{\"x\":5,\"y\":6}]}")
	unitData2 := []byte("{\"knight\":[{\"x\":10,\"y\":15},{\"x\":4,\"y\":15},{\"x\":3,\"y\":15},{\"x\":5,\"y\":15}]}")
	world.AddPlayerUnits(unitData, 0)
	world.AddPlayerUnits(unitData2, 1)

	em.Entities[2].Active = false
	em.Entities[5].Active = false
	em.Entities[6].Active = false
	em.Entities[7].Active = false

	em.Systems[0].Update()

	if len(em.Entities) != 4 {
		t.Errorf("After removing 4 entities, there should be 4 entities left, but instead got %v", len(em.Entities))
	}

}
