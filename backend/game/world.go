package game

import (
	"encoding/json"
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type World struct {
	Players       []engine.PlayerData
	EntityManager engine.EntityManagerI
	Grid          engine.Grid
	ObjectPool    *engine.ObjectPool
	EventManager  *engine.EventManager
	Tick          int
	MatchActive   bool
	//UnitDataMap        map[string]map[string]interface{}
	//EffectDataMap      map[string]map[string]interface{}
	//AbilityDataMap     map[string]map[string]interface{}
	ClientEventManager *engine.ClientEventManager
	WorkingMap         *engine.Imap
	SpatialHash        *engine.SpatialHash
	Buff               []int
	Blackboard         map[int][]int // Who targets who
}

func (w *World) World() {}

func CreateWorld() *World {
	println("World created")
	world := World{}
	world.Players = make([]engine.PlayerData, 0, 2)
	world.Tick = 0
	world.ObjectPool = engine.CreateObjectPool(30)
	world.MatchActive = true
	world.EventManager = engine.CreateEventManager()
	world.ClientEventManager = engine.CreateClientEventManager()
	world.SpatialHash = engine.CreateSpatialHash(constants.MapWidth, constants.MapHeight, math.Vector{X: float32(constants.QuadrantSize),
		Y: float32(constants.QuadrantSize)})
	world.Blackboard = make(map[int][]int, 30)

	return &world
}

func (w *World) SetupECS(e engine.EntityManagerI) {
	w.EntityManager = e
	w.EntityManager.SetObjectPool(w.ObjectPool)
	w.EntityManager.SetEventManager(w.EventManager)
}

func (w *World) GetObjectPool() *engine.ObjectPool {
	return w.ObjectPool
}

func (w *World) GetEntityManager() engine.EntityManagerI {
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

	w.Grid.Update()

	w.EntityManager.Update()

	w.checkForMatchEnd()

	w.Tick++
}

func (w *World) checkForMatchEnd() {
	for _, player := range w.Players {
		if player.NumberOfUnits == 0 {
			w.MatchActive = false
			println("Match over")
		}
	}
}

func (w *World) StartMatch() {
	w.Grid.Update()
	w.EntityManager.StartMatch()
	//mc := w.ObjectPool.Components["MovementComponent"]
	pc := w.ObjectPool.Components["PositionComponent"]

	for _, ent := range w.EntityManager.GetEntities() {
		posComp := pc[ent.Index].(components.PositionComponent)
		posComp.Address = w.SpatialHash.Insert(posComp.Address, posComp.Position, ent.ID, ent.PlayerTag)
		pc[ent.Index] = posComp
	}
}

func (w *World) AddPlayerUnits(unitData map[string][]math.Vector, tag int) {
	//Todo check if place is taken already
	for id, positions := range unitData {
		fmt.Printf("Id %s, has %v\n", id, len(unitData[id]))
		for _, pos := range positions {
			entityData := engine.NewEntityData{Data: constants.UnitDataMap[id], ID: id, PlayerTag: tag, Position: pos}
			w.EntityManager.AddEntity(entityData, tag, true)
			w.Players[tag].NumberOfUnits++
		}
	}
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
