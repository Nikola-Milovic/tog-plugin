package systems

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type MovementSystem struct {
	World *game.World
}

func (ms MovementSystem) Update() {
	world := ms.World
	//useless := 0
	//g := World.Grid
	//indexMap := world.GetEntityManager().GetIndexMap()
	entities := world.EntityManager.GetEntities()
	movementComponents := world.ObjectPool.Components["MovementComponent"]
	positionComponents := world.ObjectPool.Components["PositionComponent"]

	for index := range entities {
		if !entities[index].Active {
			continue
		}

		if entities[index].State != constants.StateWalking && entities[index].State != constants.StateEngaging {
			continue
		}

		//Setup
		playerTag := entities[index].PlayerTag

		movementComp := movementComponents[index].(components.MovementComponent)
		posComp := positionComponents[index].(components.PositionComponent)
		velocity := movementComp.Velocity

		targetPos := movementComp.Goal



	//	desiredVelocity := math.Zero()

		if movementComp.DestinationMultiplier != 0.0 {
			velocity = velocity.Add(posComp.Position.To(targetPos).Normalize().MultiplyScalar(movementComp.DestinationMultiplier))
		}

		movementComp.DestinationMultiplier += 0.2

		//Arriving
		distanceToTarget := posComp.Position.Distance(targetPos)
		arrivingZone := posComp.Radius + 80
		//
		//if distanceToTarget < arrivingZone {
		//	desiredVelocity =
		//}

		diff := distanceToTarget / arrivingZone
		if diff < 0.4 * 1/movementComp.MovementSpeed{
			diff = 0.4 * 1/movementComp.MovementSpeed
		} else if diff > 1.0 {
			diff = 1.0
		}
		velocity = velocity.Normalize().MultiplyScalar(movementComp.MovementSpeed * diff)

	//	fmt.Printf("Velocity for %s is %v\n", entities[index].UnitID, velocity)

		if !checkIfUnitInsideMap(posComp.Position.Add(velocity), posComp.Radius/2 - 2) {
			velocity = math.Zero()
			entities[index].State = constants.StateThinking
		}

		posComp.Position = posComp.Position.Add(velocity)

		data := make(map[string]interface{}, 3)
		data["event"] = "walk"
		data["who"] = entities[index].ID
		data["where"] = posComp.Position
		world.ClientEventManager.AddEvent(data)

		movementComp.Velocity = velocity

		posComp.Address = world.SpatialHash.Update(posComp.Address, posComp.Position, entities[index].ID, playerTag)

		world.ObjectPool.Components["PositionComponent"][index] = posComp
		world.ObjectPool.Components["MovementComponent"][index] = movementComp
	}

	//fmt.Println("we made ", useless, " iterations in last second")
}

func checkIfUnitInsideMap(pos math.Vector, radius float32) bool {
	isInsideMap := pos.X+radius < float32(constants.MapWidth) && pos.X-radius >= 0 && pos.Y+radius < float32(constants.MapHeight) &&
		pos.Y-radius >= 0

	return isInsideMap
}

