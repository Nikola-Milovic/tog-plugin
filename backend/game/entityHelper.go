package game

//GetNearbyEntities returns indexes of entities that are in range of maxDistance, excluding self (index parameter)
func GetNearbyEntities(maxDistance int, world *World, index int) []int {
	nearbyEntities := make([]int, 0, len(world.EntityManager.Entities))

	myPos := world.ObjectPool.Components["PositionComponents"][index].(PositionComponent)

	for idx, p := range world.ObjectPool.Components["PositionComponents"] {
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
