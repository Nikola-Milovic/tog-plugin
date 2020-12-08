package game

import (
	"encoding/json"

	"github.com/Nikola-Milovic/tog-plugin/engine"
)

//GetNearbyEntities returns indexes of entities that are in range of maxDistance, excluding self (index parameter)
func GetNearbyEntities(maxDistance int, world *World, index int) []int {
	nearbyEntities := make([]int, 0, len(world.EntityManager.Entities))

	myPos := world.ObjectPool.Components["PositionComponent"][index].(PositionComponent)

	for idx, p := range world.ObjectPool.Components["PositionComponent"] {
		posComp := p.(PositionComponent)
		if idx == index {
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
	entities := make([]engine.EntityData, 0, size+1)

	for i := 0; i < size; i++ {
		pos := e.ObjectPool.Components["PositionComponent"][i].(PositionComponent)

		//	fmt.Printf("I at %v am at position %v \n", i, e.PositionComponents[i].Position)
		entities = append(entities, engine.EntityData{
			Index:    i,
			Position: pos.Position,
			State:    e.Actions[i].GetActionType(),
			Tag:      e.Entities[i].PlayerTag,
			Path:     w.ObjectPool.Components["MovementComponent"][0].(MovementComponent).Path,
		})
	}

	data, err := json.Marshal(&entities)
	return data, err
}
