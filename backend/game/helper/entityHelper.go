package helper

import (
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//GetNearbyEntities returns indexes of entities that are in range of maxDistance, excluding self (index parameter)
func GetNearbyEntities(maxDistance int, world *game.World, index int) []int {
	nearbyEntities := make([]int, 0, len(world.EntityManager.GetEntities())+1)

	myPos := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)

	for idx, p := range world.ObjectPool.Components["PositionComponent"] {
		posComp := p.(components.PositionComponent)
		dist := world.Grid.GetDistance(posComp.Position, myPos.Position)
		//	fmt.Printf("Found entity at %v, distance to %v \n", idx, dist)
		if dist <= maxDistance {
			nearbyEntities = append(nearbyEntities, idx)
		}
	}

	return nearbyEntities
}
