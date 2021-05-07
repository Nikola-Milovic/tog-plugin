package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

type DeathSystem struct {
	World           *game.World
	IndexesToRemove []int
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

		w.EntityManager.RemoveEntity(indexMap[id]) // TODO: remove from spatial hash
	}

	for indx, comp := range w.ObjectPool.Components["StatsComponent"] {
		component := comp.(components.StatsComponent)
		if component.Health <= 0 {
			entities[indx].Active = false
		}
	}
}
