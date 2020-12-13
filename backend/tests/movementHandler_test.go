package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

// //Path from 0,0 to 0,5 should be 0,1 0,2 0,3 0,4 0,5 and after 1 movement phase should remove the 0,1
// func TestCorrectPathAndOneStep(t *testing.T) {
// 	jsonData, _ := ioutil.ReadFile("../resources/test/singleUnitTest.json")
// 	var data map[string]interface{}
// 	err := json.Unmarshal(jsonData, &data)
// 	if err != nil {
// 		t.Errorf("Couldn't unmarshal json: %e", err)
// 	}

// 	world := game.CreateWorld()

// 	world.EntityManager.AddEntity(data)
// 	world.EntityManager.AddEntity(data)

// 	//Setup 1 step TODO change later when we create a different system for movement
// 	world.Tick = constants.MovementSpeedFast

// 	//Path from 0,0 to 0,5
// 	pathToMatch := []engine.Vector{engine.Vector{0, 1}, engine.Vector{0, 2}, engine.Vector{0, 3}, engine.Vector{0, 4}}
// 	movementAction := game.MovementAction{Target: 1, Index: 0}
// 	//Setup the entities position
// 	pcomp := world.ObjectPool.Components["PositionComponent"][0].(game.PositionComponent)
// 	pcomp.Position = engine.Vector{0, 0}
// 	world.ObjectPool.Components["PositionComponent"][0] = pcomp

// 	//handle action
// 	world.EntityManager.Handlers["movement"].HandleAction(movementAction)

// 	// Result
// 	newPos := world.ObjectPool.Components["PositionComponent"][0].(game.PositionComponent)
// 	if newPos.Position.X != 0 || newPos.Position.Y != 1 {
// 		t.Errorf("After 1 step from 0,0 to 0,5 entitiy should be at X: 0, Y: 1, instead got %v", newPos.Position)
// 	}

// 	newMovementComp := world.ObjectPool.Components["MovementComponent"][0].(game.MovementComponent)
// 	if len(pathToMatch) != len(newMovementComp.Path) {
// 		t.Errorf("Path lengths are not equal, expected %v, got %v", len(pathToMatch), len(newMovementComp.Path))
// 	}
// }

func TestApproachingEachOther(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/test/twoUnitsTest.json")
	var data []interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	w := game.CreateWorld()

	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])

	w.EntityManager.Entities[0].PlayerTag = 1
	w.EntityManager.Entities[1].PlayerTag = 0

	p1 := w.ObjectPool.Components["PositionComponent"][0].(game.PositionComponent)
	p2 := w.ObjectPool.Components["PositionComponent"][1].(game.PositionComponent)

	// p1.Position = engine.Vector{0, 13}
	// p2.Position = engine.Vector{15, 13}

	p1.Position = engine.Vector{0, 7}
	p2.Position = engine.Vector{5, 7}

	w.ObjectPool.Components["PositionComponent"][0] = p1
	w.ObjectPool.Components["PositionComponent"][1] = p2
	for i := 0; i < 80; i++ {
		w.Update()
	}
}

func TestEightEntities(t *testing.T) {
	path := "../resources/units.json"

	jsonData, _ := ioutil.ReadFile(path)
	var data []interface{}

	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldnt unmarshal data %v", err.Error())
	}

	w := game.CreateWorld()
	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])

	w.EntityManager.Entities[0].PlayerTag = 1
	w.EntityManager.Entities[1].PlayerTag = 1
	w.EntityManager.Entities[2].PlayerTag = 0
	w.EntityManager.Entities[3].PlayerTag = 0
	w.EntityManager.Entities[4].PlayerTag = 1
	w.EntityManager.Entities[5].PlayerTag = 1
	w.EntityManager.Entities[6].PlayerTag = 0
	w.EntityManager.Entities[7].PlayerTag = 0

	w.AddPlayer()
	w.AddPlayer()

	w.Players[0].NumberOfUnits = 4
	w.Players[1].NumberOfUnits = 4

	p1 := w.ObjectPool.Components["PositionComponent"][0].(game.PositionComponent)
	p2 := w.ObjectPool.Components["PositionComponent"][1].(game.PositionComponent)
	p3 := w.ObjectPool.Components["PositionComponent"][2].(game.PositionComponent)
	p4 := w.ObjectPool.Components["PositionComponent"][3].(game.PositionComponent)
	p5 := w.ObjectPool.Components["PositionComponent"][4].(game.PositionComponent)
	p6 := w.ObjectPool.Components["PositionComponent"][5].(game.PositionComponent)
	p7 := w.ObjectPool.Components["PositionComponent"][6].(game.PositionComponent)
	p8 := w.ObjectPool.Components["PositionComponent"][7].(game.PositionComponent)

	p1.Position = engine.Vector{0, 0}
	p2.Position = engine.Vector{0, 3}
	p5.Position = engine.Vector{0, 7}
	p6.Position = engine.Vector{0, 13}
	p4.Position = engine.Vector{15, 1}
	p7.Position = engine.Vector{15, 4}
	p8.Position = engine.Vector{15, 6}
	p3.Position = engine.Vector{15, 13}

	w.ObjectPool.Components["PositionComponent"][0] = p1
	w.ObjectPool.Components["PositionComponent"][1] = p2
	w.ObjectPool.Components["PositionComponent"][2] = p3
	w.ObjectPool.Components["PositionComponent"][3] = p4
	w.ObjectPool.Components["PositionComponent"][4] = p5
	w.ObjectPool.Components["PositionComponent"][5] = p6
	w.ObjectPool.Components["PositionComponent"][6] = p7
	w.ObjectPool.Components["PositionComponent"][7] = p8

	h1 := w.ObjectPool.Components["HealthComponent"][1].(game.HealthComponent)
	h1.Health = 30

	h5 := w.ObjectPool.Components["HealthComponent"][5].(game.HealthComponent)
	h5.Health = 18

	h7 := w.ObjectPool.Components["HealthComponent"][7].(game.HealthComponent)
	h7.Health = 45

	w.ObjectPool.Components["HealthComponent"][1] = h1
	w.ObjectPool.Components["HealthComponent"][5] = h5
	w.ObjectPool.Components["HealthComponent"][7] = h7

	for i := 0; i < 80; i++ {
		w.Update()
	}
}
