package game

import (
	"encoding/json"
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/startup"
)

type World struct {
	Players       []engine.PlayerData
	EntityManager *engine.EntityManager
	Grid          *engine.Grid
	ObjectPool    *engine.ObjectPool
	EventManager  *engine.EventManager
	Tick          int
	MatchActive   bool
}

func CreateWorld() *World {
	println("World created")
	world := World{}
	world.Players = make([]engine.PlayerData, 0, 2)
	world.EntityManager = engine.CreateEntityManager(15)
	world.Grid = engine.CreateGrid()
	world.Tick = 0
	world.ObjectPool = engine.CreateObjectPool(15)
	world.MatchActive = true
	world.EventManager = engine.CreateEventManager()

	world.EntityManager.ObjectPool = world.ObjectPool
	world.EntityManager.EventManager = world.EventManager

	world.registerComponentMakers()
	world.registerHandlers()
	world.registerAIMakers()
	return &world
}

func (w *World) AddPlayer() int {
	tag := len(w.Players)
	w.Players = append(w.Players, engine.PlayerData{Tag: tag})
	fmt.Printf("Added player and his tag is %v\n", tag)
	return tag
}

func (w *World) registerComponentMakers() {
	w.EntityManager.RegisterComponentMaker("MovementComponent", components.MovementComponentMaker)
	w.EntityManager.RegisterComponentMaker("PositionComponent", components.PositionComponentMaker)
	w.EntityManager.RegisterComponentMaker("AttackComponent", components.AttackComponentMaker)
	w.EntityManager.RegisterComponentMaker("StatsComponent", components.StatsComponentMaker)
}

func (w *World) registerHandlers() {
	w.EntityManager.RegisterHandler(constants.MovementEvent, MovementEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.AttackEvent, AttackEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.TakeDamageEvent, TakeDamageEventHandler{World: w})
}

func (w *World) registerAIMakers() {
	w.EntityManager.RegisterAIMaker("knight", func() engine.AI { return KnightAI{world: w} })
	w.EntityManager.RegisterAIMaker("archer", func() engine.AI { return KnightAI{world: w} })
}

func (w *World) Update() {
	if !w.MatchActive {
		return
	}
	w.Tick++

	w.Grid.Update()

	w.EntityManager.Update()

	checkForDeadEntities(w)
}

func (w *World) AddPlayerUnits(data []byte, tag int) {
	unitData := make(map[string][]engine.Vector)
	//Unit data is {"knight":[{"x":1,"y":5},{"x":1,"y":6},{"x":1,"y":7},{"x":1,"y":8},{"x":1,"y":9}]}
	err := json.Unmarshal(data, &unitData)
	if err != nil {
		panic(err.Error())
	}
	//Todo check if place is taken already
	for id, positions := range unitData {
		for _, pos := range positions {
			entityData := engine.NewEntityData{Data: startup.UnitDataMap[id], ID: id, PlayerTag: tag}
			index := w.EntityManager.AddEntity(entityData)
			position := w.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
			position.Position = pos
			if tag == 1 { // Used to place the other player at the other end of the screen
				position.Position.X = w.Grid.MaxWidth/w.Grid.CellSize - position.Position.X
			}
			w.ObjectPool.Components["PositionComponent"][index] = position
			w.Players[tag].NumberOfUnits++
		}
	}
}
