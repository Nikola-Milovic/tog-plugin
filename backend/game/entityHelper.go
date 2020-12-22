package game

import (
	"encoding/json"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//GetNearbyEntities returns indexes of entities that are in range of maxDistance, excluding self (index parameter)
func GetNearbyEntities(maxDistance int, world *World, index int) []int {
	nearbyEntities := make([]int, 0, len(world.EntityManager.Entities))

	myPos := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)

	for idx, p := range world.ObjectPool.Components["PositionComponent"] {
		posComp := p.(components.PositionComponent)
		if idx == index || !world.EntityManager.Entities[idx].Active {
			continue
		}
		dist := world.Grid.GetDistance(posComp.Position, myPos.Position)
		if dist <= maxDistance {
			//	fmt.Printf("Found entity at %v, distance to %v \n", idx, dist)
			nearbyEntities = append(nearbyEntities, idx)
		}
	}

	return nearbyEntities
}

//GetEntitiesData gets the data of all entities and packs them into []byte, used to send the clients necessary data to reconstruct the current state of the game
//TODO: add batching instead of sending all the data at once
func GetEntitiesData(w *World) ([]byte, error) {
	e := w.EntityManager
	size := len(w.EntityManager.Entities)
	entities := make([]engine.EntityMessage, 0, size+1)

	for i := 0; i < size; i++ {
		pos := e.ObjectPool.Components["PositionComponent"][i].(components.PositionComponent)
		state := " "

		if !e.Entities[i].Active {
			state = "dead"
		}

		entities = append(entities, engine.EntityMessage{
			Index:    i,
			Position: pos.Position,
			State:    state,
			Tag:      e.Entities[i].PlayerTag,
			//Path:     w.ObjectPool.Components["MovementComponent"][i].(MovementComponent).Path,
		})
	}

	data, err := json.Marshal(&entities)
	return data, err
}

func checkForDeadEntities(w *World) {
	for indx, comp := range w.ObjectPool.Components["StatsComponent"] {
		if !w.EntityManager.Entities[indx].Active {
			continue
		}
		component := comp.(components.StatsComponent)
		if component.Health <= 0 {
			w.EntityManager.RemoveEntity(indx)
			w.Players[w.EntityManager.Entities[indx].PlayerTag].NumberOfUnits--
		}
	}
}
