package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

//Path from 0,0 to 0,5 should be 0,1 0,2 0,3 0,4 0,5 and after 1 movement phase should remove the 0,1
func TestAi(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/test/singleUnitTest.json")
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	world.EntityManager.AddEntity(data)

	//Setup 1 step TODO change later when we create a different system for movement
	world.Counter = constants.MovementSpeedFast

	//Path from 0,0 to 0,5
	pathToMatch := []engine.Vector{engine.Vector{0, 2}, engine.Vector{0, 2}, engine.Vector{0, 3}, engine.Vector{0, 5}}
	movementAction := game.MovementAction{Target: engine.Vector{0, 5}, Index: 0}
	//Setup the entities position
	pcomp := world.ObjectPool.Components["PositionComponent"][0].(game.PositionComponent)
	pcomp.Position = engine.Vector{0, 0}
	world.ObjectPool.Components["PositionComponent"][0] = pcomp

	//handle action
	world.EntityManager.Handlers["movement"].HandleAction(movementAction)

	// Result
	newPos := world.ObjectPool.Components["PositionComponent"][0].(game.PositionComponent)
	if newPos.Position.X != 0 || newPos.Position.Y != 1 {
		t.Errorf("After 1 step from 0,0 to 0,5 entitiy should be at X: 0, Y: 1, instead got %v", newPos.Position)
	}

	newMovementComp := world.ObjectPool.Components["MovementComponent"][0].(game.MovementComponent)
	if len(pathToMatch) != len(newMovementComp.Path) {
		t.Errorf("Path lengths are not equal, expected %v, got %v", len(pathToMatch), len(newMovementComp.Path))
	}
}

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
