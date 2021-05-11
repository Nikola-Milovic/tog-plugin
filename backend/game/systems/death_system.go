package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

type DeathSystem struct {
	World *game.World
}

//Update on deathSystem, first marks dead entities as inactive, then deletes them
func (ds DeathSystem) Update() {
	w := ds.World
	em := ds.World.EntityManager
	entities := em.GetEntities()
	indexMap := em.GetIndexMap()

	w.Buff = w.Buff[:0]

	for indx, ent := range entities {
		if !ent.Active {
			w.Buff = append(w.Buff, entities[indx].ID)
		}
	}

	for _, id := range w.Buff {
		w.Players[entities[indexMap[id]].PlayerTag].NumberOfUnits--

		//Tell client that unit died
		data := make(map[string]interface{}, 2)
		data["event"] = "death"
		data["who"] = id
		ds.World.ClientEventManager.AddEvent(data)

		ds.removeEntity(id, indexMap[id])
	}

	for indx, comp := range w.ObjectPool.Components["StatsComponent"] {
		component := comp.(components.StatsComponent)
		if component.Health <= 0 {
			entities[indx].Active = false
		}
	}
}

func (ds DeathSystem) removeEntity(id, index int) {
	address := ds.World.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent).Address

	ds.World.EntityManager.RemoveEntity(index) // TODO: remove from spatial hash and blackboard

	ds.World.SpatialHash.Remove(address, id)

	idx := -1
	for key, ids := range ds.World.Blackboard {
		for idindex, i := range ids {
			if i == id {
				idx = idindex
			}
		}
		ds.World.Blackboard[key] = removeIndex(ids, idx)
	}

	delete(ds.World.Blackboard, id)

}
