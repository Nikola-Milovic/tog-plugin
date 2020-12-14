package game

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type World struct {
	Players       []engine.PlayerData
	EntityManager *engine.EntityManager
	Grid          *engine.Grid
	ObjectPool    *engine.ObjectPool
	Tick          int
	MatchActive   bool
}

func CreateWorld() *World {
	println("World created")
	world := World{}
	world.Players = make([]engine.PlayerData, 0, 2)
	world.EntityManager = engine.CreateEntityManager()
	world.Grid = engine.CreateGrid()
	world.Tick = 0
	world.ObjectPool = engine.CreateObjectPool(10)
	world.MatchActive = true

	world.EntityManager.ObjectPool = world.ObjectPool
	world.registerComponentMakers()
	world.registerHandlers()
	world.registerAIMakers()
	return &world
}

func (w *World) AddPlayer() int {
	tag := len(w.Players)
	w.Players = append(w.Players, engine.PlayerData{Tag: tag})
	return tag
}

func (w *World) registerComponentMakers() {
	w.EntityManager.RegisterComponentMaker("MovementComponent", MovementComponentMaker)
	w.EntityManager.RegisterComponentMaker("PositionComponent", PositionComponentMaker)
	w.EntityManager.RegisterComponentMaker("AttackComponent", AttackComponentMaker)
	w.EntityManager.RegisterComponentMaker("HealthComponent", HealthComponentMaker)
}

func (w *World) registerHandlers() {
	w.EntityManager.RegisterHandler(constants.ActionTypeMovement, MovementHandler{world: w})
	w.EntityManager.RegisterHandler(constants.ActionTypeAttack, AttackHandler{world: w})
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
