package game

import (
	"encoding/json"
	"fmt"

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

//GetClientEvents has
//TODO: add batching instead of sending all the data at once
func GetClientEvents(w *World) ([]byte, error) {

	events := w.ClientEventManager.Events

	data, err := json.Marshal(&events)

	w.ClientEventManager.Events = w.ClientEventManager.Events[:0]

	if err != nil {
		fmt.Printf("Error marshaling client events is %v", err.Error())
		return nil, err
	}

	return data, err
}
