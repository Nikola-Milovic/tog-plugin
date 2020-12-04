package game

import "github.com/Nikola-Milovic/tog-plugin/engine"

type World struct {
	EntityManager *engine.EntityManager
	Grid          *engine.Grid
	ObjectPool    *engine.ObjectPool
	Counter       engine.Counter
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
	return &world
}

func (w *World) registerComponentMakers() {
	w.EntityManager.RegisterComponentMaker("MovementComponent", MovementComponentMaker)
	w.EntityManager.RegisterComponentMaker("PositionComponent", PositionComponentMaker)
}

func (w *World) registerHandlers() {
	w.EntityManager.RegisterHandler("movement", MovementHandler{world: w})
	w.EntityManager.RegisterHandler("attack", AttackHandler{world: w})
}
