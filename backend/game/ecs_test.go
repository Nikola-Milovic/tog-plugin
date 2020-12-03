package game

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestEntityCreation(t *testing.T) {
	plan, _ := ioutil.ReadFile("../resources/units.json")
	var data map[string]interface{}
	err := json.Unmarshal(plan, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := CreateWorld()
	
	world.EntityManager.AddEntity(data["Knight"].(map[string]interface{}))
}

//int(data["Components"].(map[string]interface {})["AttackComponent"].(map[string]interface {})["damage"].(float64)))
