package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

func TestEightEntities(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := game.CreateWorld()

	world.AddPlayer()
	world.AddPlayer()

	unitData := []byte("{\"knight\":[{\"x\":5,\"y\":3},{\"x\":4,\"y\":6},{\"x\":3,\"y\":6},{\"x\":5,\"y\":6}]}")
	unitData2 := []byte("{\"knight\":[{\"x\":10,\"y\":15},{\"x\":4,\"y\":15},{\"x\":3,\"y\":15},{\"x\":5,\"y\":15}]}")
	world.AddPlayerUnits(unitData, 0)
	world.AddPlayerUnits(unitData2, 1)

	p1 := world.ObjectPool.Components["PositionComponent"][0].(game.PositionComponent)
	p2 := world.ObjectPool.Components["PositionComponent"][1].(game.PositionComponent)
	p3 := world.ObjectPool.Components["PositionComponent"][2].(game.PositionComponent)
	p4 := world.ObjectPool.Components["PositionComponent"][3].(game.PositionComponent)
	p5 := world.ObjectPool.Components["PositionComponent"][4].(game.PositionComponent)
	p6 := world.ObjectPool.Components["PositionComponent"][5].(game.PositionComponent)
	p7 := world.ObjectPool.Components["PositionComponent"][6].(game.PositionComponent)
	p8 := world.ObjectPool.Components["PositionComponent"][7].(game.PositionComponent)

	p1.Position = engine.Vector{0, 0}
	p2.Position = engine.Vector{0, 3}
	p5.Position = engine.Vector{0, 7}
	p6.Position = engine.Vector{0, 13}
	p4.Position = engine.Vector{15, 1}
	p7.Position = engine.Vector{15, 4}
	p8.Position = engine.Vector{15, 6}
	p3.Position = engine.Vector{15, 13}

	world.ObjectPool.Components["PositionComponent"][0] = p1
	world.ObjectPool.Components["PositionComponent"][1] = p2
	world.ObjectPool.Components["PositionComponent"][2] = p3
	world.ObjectPool.Components["PositionComponent"][3] = p4
	world.ObjectPool.Components["PositionComponent"][4] = p5
	world.ObjectPool.Components["PositionComponent"][5] = p6
	world.ObjectPool.Components["PositionComponent"][6] = p7
	world.ObjectPool.Components["PositionComponent"][7] = p8

	h1 := world.ObjectPool.Components["HealthComponent"][1].(game.HealthComponent)
	h1.Health = 30

	h5 := world.ObjectPool.Components["HealthComponent"][5].(game.HealthComponent)
	h5.Health = 18

	h7 := world.ObjectPool.Components["HealthComponent"][7].(game.HealthComponent)
	h7.Health = 45

	world.ObjectPool.Components["HealthComponent"][1] = h1
	world.ObjectPool.Components["HealthComponent"][5] = h5
	world.ObjectPool.Components["HealthComponent"][7] = h7

	for i := 0; i < 80; i++ {
		world.Update()
	}
}

// // //Path from 0,0 to 0,5 should be 0,1 0,2 0,3 0,4 0,5 and after 1 movement phase should remove the 0,1
// // func TestCorrectPathAndOneStep(t *testing.T) {
// // 	jsonData, _ := ioutil.ReadFile("../resources/test/singleUnitTest.json")
// // 	var data map[string]interface{}
// // 	err := json.Unmarshal(jsonData, &data)
// // 	if err != nil {
// // 		t.Errorf("Couldn't unmarshal json: %e", err)
// // 	}

// // 	world := game.CreateWorld()

// // 	world.EntityManager.AddEntity(data)
// // 	world.EntityManager.AddEntity(data)

// // 	//Setup 1 step TODO change later when we create a different system for movement
// // 	world.Tick = constants.MovementSpeedFast

// // 	//Path from 0,0 to 0,5
// // 	pathToMatch := []engine.Vector{engine.Vector{0, 1}, engine.Vector{0, 2}, engine.Vector{0, 3}, engine.Vector{0, 4}}
// // 	movementAction := game.MovementAction{Target: 1, Index: 0}
// // 	//Setup the entities position
// // 	pcomp := world.ObjectPool.Components["PositionComponent"][0].(game.PositionComponent)
// // 	pcomp.Position = engine.Vector{0, 0}
// // 	world.ObjectPool.Components["PositionComponent"][0] = pcomp

// // 	//handle action
// // 	world.EntityManager.Handlers["movement"].HandleAction(movementAction)

// // 	// Result
// // 	newPos := world.ObjectPool.Components["PositionComponent"][0].(game.PositionComponent)
// // 	if newPos.Position.X != 0 || newPos.Position.Y != 1 {
// // 		t.Errorf("After 1 step from 0,0 to 0,5 entitiy should be at X: 0, Y: 1, instead got %v", newPos.Position)
// // 	}

// // 	newMovementComp := world.ObjectPool.Components["MovementComponent"][0].(game.MovementComponent)
// // 	if len(pathToMatch) != len(newMovementComp.Path) {
// // 		t.Errorf("Path lengths are not equal, expected %v, got %v", len(pathToMatch), len(newMovementComp.Path))
// // 	}
// // }

// func TestApproachingEachOther(t *testing.T) {
// 	jsonData, _ := ioutil.ReadFile("../resources/test/twoUnitsTest.json")
// 	var data []interface{}
// 	err := json.Unmarshal(jsonData, &data)
// 	if err != nil {
// 		t.Errorf("Couldn't unmarshal json: %e", err)
// 	}

// 	w := game.CreateWorld()

// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])

// 	w.EntityManager.Entities[0].PlayerTag = 1
// 	w.EntityManager.Entities[1].PlayerTag = 0

// 	p1 := w.ObjectPool.Components["PositionComponent"][0].(game.PositionComponent)
// 	p2 := w.ObjectPool.Components["PositionComponent"][1].(game.PositionComponent)

// 	// p1.Position = engine.Vector{0, 13}
// 	// p2.Position = engine.Vector{15, 13}

// 	p1.Position = engine.Vector{0, 7}
// 	p2.Position = engine.Vector{5, 7}

// 	w.ObjectPool.Components["PositionComponent"][0] = p1
// 	w.ObjectPool.Components["PositionComponent"][1] = p2
// 	for i := 0; i < 80; i++ {
// 		w.Update()
// 	}
// }
