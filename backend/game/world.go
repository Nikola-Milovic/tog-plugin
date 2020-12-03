package game

import "github.com/Nikola-Milovic/tog-plugin/engine"

type World struct {
	EntityManager *engine.EntityManager
	Grid          *engine.Grid
	ObjectPool    *engine.ObjectPool
	Counter       engine.Counter
}

func CreateWorld() *World {
	world := World{}
	world.EntityManager = engine.CreateEntityManager()
	world.Grid = engine.CreateGrid()
	world.Counter = 0
	world.ObjectPool = engine.CreateObjectPool(10)

	world.EntityManager.ObjectPool = world.ObjectPool

	return &world
}
