package tempsys

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

type LineShotTempSystem struct {
	World *game.World
	Data  []map[string]interface{}
}

func (ds *LineShotTempSystem) Update() {
	if len(ds.Data) == 0 {
		ds.World.EntityManager.RemoveTempSystem("LineshotTempSystem")
		return
	}

	expired := make([]int, 0, 10)

	for indx, data := range ds.Data {
		//	fmt.Printf("data %v", data)
		abId := data["ability_id"].(string)
		speed := data["speed"].(int)
		casterID := data["caster"].(string)
		destination := data["target"].(engine.Vector)
		position := data["position"].(engine.Vector)
		lastMovedTick := data["last_moved"].(int)

		cell, _ := ds.World.Grid.CellAt(position)
		// 1) Check for collision
		if ds.World.Grid.IsCellTaken(position) && cell.OccupiedID != casterID {
			fmt.Printf("Collision with %v, at %v\n", cell.OccupiedID, position)
			ds.OnCollision(abId, cell.OccupiedID)
			expired = append(expired, indx)
			continue
		}

		// 2) Check for expiration/ end of range
		if destination == position {
			expired = append(expired, indx)
		}

		// 3) Check for movement
		if ds.World.Tick-lastMovedTick > speed {
			lastMovedTick = ds.World.Tick
			if destination.X == position.X {
				if position.Y > destination.Y {
					position.Y = position.Y - 1
				} else {
					position.Y = position.Y + 1
				}
			} else {
				if position.X > destination.X {
					position.X = position.X - 1
				} else {
					position.X = position.X + 1
				}
			}

			fmt.Printf("Projectile target is to %v\n", destination)
			fmt.Printf("Projectile moving from to %v\n", position)

			data["position"] = position
			data["last_moved"] = lastMovedTick

			// clientEvent := make(map[string]interface{}, 2)
			// clientEvent["projectile_id"] = data["id"]
			// clientEvent["event"] = data["projectile_movement"]
			// clientEvent["position"] = data["position"]

			// ds.World.ClientEventManager.AddEvent(clientEvent)
		}

	}

	for _, index := range expired {
		clientEvent := make(map[string]interface{}, 2)
		clientEvent["id"] = ds.Data[index]["id"]
		clientEvent["event"] = "projectile_expired"

		ds.Data[index] = ds.Data[len(ds.Data)-1]
		ds.Data = ds.Data[:len(ds.Data)-1]
		fmt.Printf("Expired %v\n", index)

		ds.World.ClientEventManager.AddEvent(clientEvent)
	}

}

func (ds LineShotTempSystem) OnCollision(abID string, collisionID string) {
	onHit := ds.World.AbilityDataMap[abID]["OnHit"].(map[string]interface{})

	acts := onHit["Actions"].([]interface{})
	for _, act := range acts {
		data := make(map[string]interface{}, 5)
		data["target"] = collisionID
		data["action_data"] = act
		ev := engine.Event{}
		ev.ID = constants.TriggerActionEvent
		ev.Index = -1
		ev.Priority = constants.TriggerActionEventPriority
		ev.Data = data

		ds.World.EventManager.SendEvent(ev)

	}
}

func (ds *LineShotTempSystem) AddData(data map[string]interface{}) {
	ds.Data = append(ds.Data, data)
}

func CreateLineShotTempSystem(w interface{}) engine.TempSystem {
	world := w.(*game.World)
	fmt.Print(world)
	return &LineShotTempSystem{World: world, Data: make([]map[string]interface{}, 0, 10)}
}
