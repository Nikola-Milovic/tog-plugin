package engine

type World struct {
	EntityManager *EntityManager
	Grid          *Grid
	ObjectPool    *ObjectPool
	Counter       Counter
}
