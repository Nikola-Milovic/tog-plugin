package game

import (
	"encoding/json"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//GetNearbyEntities returns indexes of entities that are in range of maxDistance, excluding self (index parameter)
func GetNearbyEntities(maxDistance int, world *World, index int) []int {
	nearbyEntities := make([]int, 0, len(world.EntityManager.Entities)+1)

	myPos := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)

	for idx, p := range world.ObjectPool.Components["PositionComponent"] {
		posComp := p.(components.PositionComponent)
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

	for i, ent := range w.EntityManager.Entities {
		pos := e.ObjectPool.Components["PositionComponent"][i].(components.PositionComponent)
		state := " "

		if !ent.Active {
			state = "dead"
		}

		entities = append(entities, engine.EntityMessage{
			ID:       ent.ID,
			Position: pos.Position,
			State:    state,
			//Path:     w.ObjectPool.Components["MovementComponent"][i].(MovementComponent).Path,
			Health: w.ObjectPool.Components["StatsComponent"][i].(components.StatsComponent).Health,
		})
	}

	data, err := json.Marshal(&entities)
	return data, err
}
