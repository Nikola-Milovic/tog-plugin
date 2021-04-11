//
// if atkComp.IsAttacking && w.Tick-atkComp.TimeSinceLastAttack >= atkComp.AttackSpeed {
// 	atkComp.IsAttacking = false
// }

// if movComp.IsMoving && w.Tick-movComp.TimeSinceLastMovement >= movComp.MovementSpeed {
// 	movComp.IsMoving = false
// }
package systems

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

type MovementSystem struct {
	World *game.World
}

func (ms MovementSystem) Update() {
	world := ms.World
	for index, comp := range world.ObjectPool.Components["MovementComponent"] {
		movementComp := comp.(components.MovementComponent)
		if !movementComp.IsMoving {
			continue
		}

		if !(ms.World.Tick-movementComp.TimeSinceLastMovement > movementComp.MovementSpeed) {
			continue
		}

		entities := world.EntityManager.GetEntities()
		//Finished moving

		positionComp := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)

		path := movementComp.Path

		posToMove := path[len(path)-1]

		// if world.Grid.IsCellTaken(posToMove) { //Our position is taken, we cannot go further, reset the movement
		// 	movementComp.IsMoving = false
		// 	movementComp.TimeSinceLastMovement = 0
		// 	movementComp.Path = movementComp.Path[:0]
		// 	world.ObjectPool.Components["MovementComponent"][index] = movementComp
		// 	continue
		// }

		world.Grid.ReleaseCell(positionComp.Position)

		positionComp.Position = posToMove

		world.Grid.OccupyCell(posToMove, entities[index].ID)

		path = path[:len(path)-1]

		movementComp.Path = path
		movementComp.IsMoving = false

		fmt.Printf("I %v am at %v on tick %v\n", entities[index].ID, posToMove, world.Tick)

		world.ObjectPool.Components["MovementComponent"][index] = movementComp
		world.ObjectPool.Components["PositionComponent"][index] = positionComp
	}
}
