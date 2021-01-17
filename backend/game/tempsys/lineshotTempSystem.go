package tempsys

import (
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/game"
)

type LineShotTempSystem struct {
	World *game.World
	Data  []map[string]interface{}
}

func (ds *LineShotTempSystem) Update() {
	if len(ds.Data) == 0 {
		ds.World.EntityManager.RemoveTempSystem("LineShotTempSystem")
		return
	}

	for _, data := range ds.Data {
		fmt.Printf("data %v", data)
		// 1) Check for collision

		// 2) Check for expiration/ end of range

		// 3) Check for movement

	}

}

func (ds *LineShotTempSystem) AddData(data map[string]interface{}) {
	ds.Data = append(ds.Data, data)
}
