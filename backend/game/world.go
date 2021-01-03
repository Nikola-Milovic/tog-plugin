package game

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/startup"
)

type World struct {
	Players        []engine.PlayerData
	EntityManager  *engine.EntityManager
	Grid           *engine.Grid
	ObjectPool     *engine.ObjectPool
	EventManager   *engine.EventManager
	Tick           int
	MatchActive    bool
	UnitDataMap    map[string]map[string]interface{}
	EffectDataMap  map[string]map[string]interface{}
	AbilityDataMap map[string]map[string]interface{}
}

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

	world.EntityManager.ObjectPool = world.ObjectPool
	world.EntityManager.EventManager = world.EventManager

	//Copy the data maps from startup so each match accesses its own data
	world.EffectDataMap = engine.CopyJsonMap(startup.EffectDataMap)
	world.UnitDataMap = engine.CopyJsonMap(startup.UnitDataMap)
	world.AbilityDataMap = engine.CopyJsonMap(startup.AbilityDataMap)

	return &world
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
			entityData := engine.NewEntityData{Data: w.UnitDataMap[id], ID: id, PlayerTag: tag}
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
