package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestTwoEntitiesApproaching(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5}]}}")

	world := CreateTestWorld(u1, u2, t)

	for world.MatchActive {
		world.Update()
	}
}
