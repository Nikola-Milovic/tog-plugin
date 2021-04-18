package helper

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

//GetNearbyEntities returns indexes of entities that are in range of maxDistance, excluding self (index parameter)
func GetNearbyEntities(maxDistance float32, world *game.World, index int) []int {
	nearbyEntities := make([]int, 0, len(world.EntityManager.GetEntities())+1)

	myPos := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)

	for idx, p := range world.ObjectPool.Components["PositionComponent"] {
		if idx == index {
			continue
		}
		posComp := p.(components.PositionComponent)
		dist := engine.GetDistanceIncludingDiagonalVectors(posComp.Position, myPos.Position)
		//	fmt.Printf("Found entity at %v, distance to %v \n", idx, dist)
		if dist <= maxDistance {
			nearbyEntities = append(nearbyEntities, idx)
		}
	}

	return nearbyEntities
}

//GetTaggedNearbyEntities returns indexes of entities that are in range of maxDistance, excluding self (index parameter)
func GetTaggedNearbyEntities(maxDistance float32, world *game.World, index int, tag int) []int {
	nearbyEntities := make([]int, 0, len(world.EntityManager.GetEntities())+1)

	myPos := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	positions := world.ObjectPool.Components["PositionComponent"]

	for idx, ent := range world.GetEntityManager().GetEntities(){
		if idx == index || ent.PlayerTag != tag{
			continue
		}
		posComp := positions[ent.Index].(components.PositionComponent)
		dist := engine.GetDistanceIncludingDiagonalVectors(posComp.Position, myPos.Position)
		//	fmt.Printf("Found entity at %v, distance to %v \n", idx, dist)
		if dist <= maxDistance {
			nearbyEntities = append(nearbyEntities, idx)
		}
	}

	return nearbyEntities
}

