// package game

// import (
// 	"encoding/json"
// 	"io/ioutil"

// 	"github.com/Nikola-Milovic/tog-plugin/engine"
// 	"github.com/heroiclabs/nakama-common/runtime"
// )

// func testGame(w *World, logger runtime.Logger) {

// 	path := "/nakama/data/units.json"

// 	jsonData, _ := ioutil.ReadFile(path)
// 	var data []interface{}

// 	err := json.Unmarshal(jsonData, &data)
// 	if err != nil {
// 		logger.Error("Couldn't unmarshal json: %e", err.Error())
// 		return
// 	}
// 	logger.Debug("Unit data is %v", data)

// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])

// 	w.EntityManager.Entities[0].PlayerTag = 1
// 	w.EntityManager.Entities[1].PlayerTag = 1
// 	w.EntityManager.Entities[2].PlayerTag = 0
// 	w.EntityManager.Entities[3].PlayerTag = 0
// 	w.EntityManager.Entities[4].PlayerTag = 1
// 	w.EntityManager.Entities[5].PlayerTag = 1
// 	w.EntityManager.Entities[6].PlayerTag = 0
// 	w.EntityManager.Entities[7].PlayerTag = 0

// 	w.AddPlayer()
// 	w.AddPlayer()

// 	w.Players[0].NumberOfUnits = 4
// 	w.Players[1].NumberOfUnits = 4

// 	p1 := w.ObjectPool.Components["PositionComponent"][0].(PositionComponent)
// 	p2 := w.ObjectPool.Components["PositionComponent"][1].(PositionComponent)
// 	p3 := w.ObjectPool.Components["PositionComponent"][2].(PositionComponent)
// 	p4 := w.ObjectPool.Components["PositionComponent"][3].(PositionComponent)
// 	p5 := w.ObjectPool.Components["PositionComponent"][4].(PositionComponent)
// 	p6 := w.ObjectPool.Components["PositionComponent"][5].(PositionComponent)
// 	p7 := w.ObjectPool.Components["PositionComponent"][6].(PositionComponent)
// 	p8 := w.ObjectPool.Components["PositionComponent"][7].(PositionComponent)

// 	p1.Position = engine.Vector{0, 0}
// 	p2.Position = engine.Vector{0, 3}
// 	p5.Position = engine.Vector{0, 7}
// 	p6.Position = engine.Vector{0, 13}
// 	p4.Position = engine.Vector{15, 1}
// 	p7.Position = engine.Vector{15, 4}
// 	p8.Position = engine.Vector{15, 6}
// 	p3.Position = engine.Vector{15, 13}

// 	w.ObjectPool.Components["PositionComponent"][0] = p1
// 	w.ObjectPool.Components["PositionComponent"][1] = p2
// 	w.ObjectPool.Components["PositionComponent"][2] = p3
// 	w.ObjectPool.Components["PositionComponent"][3] = p4
// 	w.ObjectPool.Components["PositionComponent"][4] = p5
// 	w.ObjectPool.Components["PositionComponent"][5] = p6
// 	w.ObjectPool.Components["PositionComponent"][6] = p7
// 	w.ObjectPool.Components["PositionComponent"][7] = p8

// 	h1 := w.ObjectPool.Components["HealthComponent"][1].(HealthComponent)
// 	h1.Health = 30

// 	h5 := w.ObjectPool.Components["HealthComponent"][5].(HealthComponent)
// 	h5.Health = 18

// 	h7 := w.ObjectPool.Components["HealthComponent"][7].(HealthComponent)
// 	h7.Health = 45

// 	w.ObjectPool.Components["HealthComponent"][1] = h1
// 	w.ObjectPool.Components["HealthComponent"][5] = h5
// 	w.ObjectPool.Components["HealthComponent"][7] = h7

// }

// func testGame2(w *World, logger runtime.Logger) {
// 	path := "/nakama/data/units.json"

// 	jsonData, _ := ioutil.ReadFile(path)
// 	var data []interface{}

// 	err := json.Unmarshal(jsonData, &data)
// 	if err != nil {
// 		logger.Error("Couldn't unmarshal json: %e", err.Error())
// 		return
// 	}
// 	logger.Debug("Unit data is %v", data)

// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])

// 	w.EntityManager.Entities[0].PlayerTag = 1
// 	w.EntityManager.Entities[1].PlayerTag = 0
// 	w.EntityManager.Entities[2].PlayerTag = 0
// 	w.EntityManager.Entities[3].PlayerTag = 0
// 	w.EntityManager.Entities[4].PlayerTag = 0
// 	w.EntityManager.Entities[5].PlayerTag = 0
// 	w.EntityManager.Entities[6].PlayerTag = 0
// 	w.EntityManager.Entities[7].PlayerTag = 0

// 	p1 := w.ObjectPool.Components["PositionComponent"][0].(PositionComponent)
// 	p2 := w.ObjectPool.Components["PositionComponent"][1].(PositionComponent)
// 	p3 := w.ObjectPool.Components["PositionComponent"][2].(PositionComponent)
// 	p4 := w.ObjectPool.Components["PositionComponent"][3].(PositionComponent)
// 	p5 := w.ObjectPool.Components["PositionComponent"][4].(PositionComponent)
// 	p6 := w.ObjectPool.Components["PositionComponent"][5].(PositionComponent)
// 	p7 := w.ObjectPool.Components["PositionComponent"][6].(PositionComponent)
// 	p8 := w.ObjectPool.Components["PositionComponent"][7].(PositionComponent)

// 	p1.Position = engine.Vector{10, 10}
// 	p2.Position = engine.Vector{3, 3}
// 	p5.Position = engine.Vector{5, 7}
// 	p6.Position = engine.Vector{0, 13}
// 	p4.Position = engine.Vector{15, 1}
// 	p7.Position = engine.Vector{15, 4}
// 	p8.Position = engine.Vector{15, 6}
// 	p3.Position = engine.Vector{15, 13}

// 	w.ObjectPool.Components["PositionComponent"][0] = p1
// 	w.ObjectPool.Components["PositionComponent"][1] = p2
// 	w.ObjectPool.Components["PositionComponent"][2] = p3
// 	w.ObjectPool.Components["PositionComponent"][3] = p4
// 	w.ObjectPool.Components["PositionComponent"][4] = p5
// 	w.ObjectPool.Components["PositionComponent"][5] = p6
// 	w.ObjectPool.Components["PositionComponent"][6] = p7
// 	w.ObjectPool.Components["PositionComponent"][7] = p8
// }

// func testGame3(w *World, logger runtime.Logger) {
// 	path := "/nakama/data/units.json"

// 	jsonData, _ := ioutil.ReadFile(path)
// 	var data []interface{}

// 	err := json.Unmarshal(jsonData, &data)
// 	if err != nil {
// 		logger.Error("Couldn't unmarshal json: %e", err.Error())
// 		return
// 	}
// 	logger.Debug("Unit data is %v", data)

// 	w.EntityManager.AddEntity(data[0])
// 	w.EntityManager.AddEntity(data[0])

// 	w.EntityManager.Entities[0].PlayerTag = 1
// 	w.EntityManager.Entities[1].PlayerTag = 0

// 	p1 := w.ObjectPool.Components["PositionComponent"][0].(PositionComponent)
// 	p2 := w.ObjectPool.Components["PositionComponent"][1].(PositionComponent)

// 	w.AddPlayer()
// 	w.AddPlayer()

// 	w.Players[0].NumberOfUnits = 1
// 	w.Players[1].NumberOfUnits = 1

// 	p1.Position = engine.Vector{0, 5}
// 	p2.Position = engine.Vector{5, 5}

// 	w.ObjectPool.Components["PositionComponent"][0] = p1
// 	w.ObjectPool.Components["PositionComponent"][1] = p2

// 	h1 := w.ObjectPool.Components["HealthComponent"][1].(HealthComponent)
// 	h1.Health = 30
// 	w.ObjectPool.Components["HealthComponent"][1] = h1

// }

package game
