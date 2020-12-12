package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/Nikola-Milovic/tog-plugin/game"
)

func TestTime(t *testing.T) {
	start := time.Now()

	jsonData, _ := ioutil.ReadFile("../resources/test/singleUnitTest.json")
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	for i := 0; i < 5; i++ {
		world.EntityManager.AddEntity(data)
		world.EntityManager.Entities[i].PlayerTag = byte(i % 2)
	}

	end := time.Now()

	fmt.Printf("Total time to execute is %v", end.Sub(start))
}
