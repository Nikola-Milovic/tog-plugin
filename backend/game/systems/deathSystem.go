package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

type DeathSystem struct {
	World           *game.World
	IndexesToRemove []string
}

//Update on deathSystem, first marks dead entities as inactive, then deletes them
func (ds DeathSystem) Update() {
	w := ds.World
	em := ds.World.EntityManager

	for indx, ent := range em.Entities {
		if !ent.Active {
			ds.IndexesToRemove = append(ds.IndexesToRemove, em.Entities[indx].ID)
		}
	}

	for _, id := range ds.IndexesToRemove {
		w.Players[em.Entities[em.IndexMap[id]].PlayerTag].NumberOfUnits--

		//Tell client that unit died
		data := make(map[string]interface{}, 2)
		data["event"] = "death"
		data["who"] = id
		ds.World.ClientEventManager.AddEvent(data)

		w.EntityManager.RemoveEntity(em.IndexMap[id])
	}

	for indx, comp := range w.ObjectPool.Components["StatsComponent"] {
		component := comp.(components.StatsComponent)
		if component.Health <= 0 {
			em.Entities[indx].Active = false
		}
	}

	ds.IndexesToRemove = ds.IndexesToRemove[:0]
}
