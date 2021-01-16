package game

import (
	"encoding/json"
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/startup"
)

type World struct {
	Players            []engine.PlayerData
	EntityManager      *engine.EntityManager
	Grid               *engine.Grid
	ObjectPool         *engine.ObjectPool
	EventManager       *engine.EventManager
	Tick               int
	MatchActive        bool
	UnitDataMap        map[string]map[string]interface{}
	EffectDataMap      map[string]map[string]interface{}
	AbilityDataMap     map[string]map[string]interface{}
	ClientEventManager *engine.ClientEventManager
}

func (w *World) World() {}

func CreateWorld() *World {
	println("World created")
	world := World{}
	world.Players = make([]engine.PlayerData, 0, 2)
	world.EntityManager = engine.CreateEntityManager(30)
	world.Grid = engine.CreateGrid()
	world.Tick = 0
	world.ObjectPool = engine.CreateObjectPool(30)
	world.MatchActive = true
	world.EventManager = engine.CreateEventManager()
	world.ClientEventManager = engine.CreateClientEventManager()

	world.EntityManager.ObjectPool = world.ObjectPool
	world.EntityManager.EventManager = world.EventManager

	//Copy the data maps from startup so each match accesses its own data
	world.EffectDataMap = engine.CopyJsonMap(startup.EffectDataMap)
	world.UnitDataMap = engine.CopyJsonMap(startup.UnitDataMap)
	world.AbilityDataMap = engine.CopyJsonMap(startup.AbilityDataMap)

	return &world
}

func (w *World) GetObjectPool() *engine.ObjectPool {
	return w.ObjectPool
}

func (w *World) GetEntityManager() *engine.EntityManager {
	return w.EntityManager
}

func (w *World) AddPlayer(id string) int {
	tag := len(w.Players)
	w.Players = append(w.Players, engine.PlayerData{Tag: tag, ID: id})
	fmt.Printf("Added player and his tag is %v\n", tag)
	return tag
}

func (w *World) Update() {
	if !w.MatchActive {
		return
	}
	w.Tick++

	w.Grid.Update()

	w.EntityManager.Update()

	w.checkForMatchEnd()
}

func (w *World) checkForMatchEnd() {
	for _, player := range w.Players {
		if player.NumberOfUnits == 0 {
			w.MatchActive = false
			println("Match over")
		}
	}
}

func (w *World) AddPlayerUnits(unitData map[string][]engine.Vector, tag int) {
	//Todo check if place is taken already
	for id, positions := range unitData {
		fmt.Printf("Id %s, has %v\n", id, len(unitData[id]))
		for _, pos := range positions {
			entityData := engine.NewEntityData{Data: w.UnitDataMap[id], ID: id, PlayerTag: tag, Position: pos}
			w.EntityManager.AddEntity(entityData, tag, true)
			w.Players[tag].NumberOfUnits++
		}
	}
}

func (w *World) GetUnitDataMap() map[string]map[string]interface{} {
	return w.UnitDataMap
}
func (w *World) GetAbilityDataMap() map[string]map[string]interface{} {
	return w.AbilityDataMap
}
func (w *World) GetEffectDataMap() map[string]map[string]interface{} {
	return w.EffectDataMap
}

//GetClientEvents has
//TODO: add batching instead of sending all the data at once
func (w *World) GetClientEvents() ([]byte, error) {
	events := w.ClientEventManager.Events

	data, err := json.Marshal(&events)

	w.ClientEventManager.Events = w.ClientEventManager.Events[:0]

	if err != nil {
		fmt.Printf("Error marshaling client events is %v", err.Error())
		return nil, err
	}

	return data, err
}
