//
// if atkComp.IsAttacking && w.Tick-atkComp.TimeSinceLastAttack >= atkComp.AttackSpeed {
// 	atkComp.IsAttacking = false
// }

// if movComp.IsMoving && w.Tick-movComp.TimeSinceLastMovement >= movComp.MovementSpeed {
// 	movComp.IsMoving = false
// }
package systems

import (
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

		//Finished moving

		positionComp := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)

		path := movementComp.Path

		world.Grid.ReleaseCell(positionComp.Position)

		posToMove := path[len(path)-1]

		if world.Grid.IsCellTaken(posToMove) {
			movementComp.IsMoving = false
			world.ObjectPool.Components["MovementComponent"][index] = movementComp
			continue
		}

		positionComp.Position = posToMove

		world.Grid.OccupyCell(posToMove)

		path = path[:len(path)-1]

		movementComp.Path = path
		movementComp.IsMoving = false

		world.ObjectPool.Components["MovementComponent"][index] = movementComp
		world.ObjectPool.Components["PositionComponent"][index] = positionComp
	}
}