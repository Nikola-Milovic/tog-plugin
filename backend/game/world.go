package game

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

type World struct {
	EntityManager *engine.EntityManager
	Grid          *engine.Grid
	ObjectPool    *engine.ObjectPool
	Counter       int
}

func CreateWorld() *World {
	println("World created")
	world := World{}
	world.EntityManager = engine.CreateEntityManager()
	world.Grid = engine.CreateGrid()
	world.Counter = 0
	world.ObjectPool = engine.CreateObjectPool(10)

	world.EntityManager.ObjectPool = world.ObjectPool
	world.registerComponentMakers()
	world.registerHandlers()
	world.registerAIMakers()
	return &world
}

func (w *World) registerComponentMakers() {
	w.EntityManager.RegisterComponentMaker("MovementComponent", MovementComponentMaker)
	w.EntityManager.RegisterComponentMaker("PositionComponent", PositionComponentMaker)
	w.EntityManager.RegisterComponentMaker("AttackComponent", AttackComponentMaker)
}

func (w *World) registerHandlers() {
	w.EntityManager.RegisterHandler(constants.ActionTypeMovement, MovementHandler{world: w})
	w.EntityManager.RegisterHandler(constants.ActionTypeAttack, AttackHandler{world: w})
}

func (w *World) registerAIMakers() {
	w.EntityManager.RegisterAIMaker("knight", func() engine.AI { return KnightAI{world: w} })
	w.EntityManager.RegisterAIMaker("archer", func() engine.AI { return KnightAI{world: w} })
}
