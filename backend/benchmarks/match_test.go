package benchmarks

import (
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/match"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
	"testing"
)

func CreateTestWorld(unitData []byte, unitData2 []byte) *game.World {
	world := game.CreateWorld()
	registry.RegisterWorld(world)

	world.AddPlayer("")
	world.AddPlayer("")

	data1 := match.PlayerReadyDataMessage{}
	data2 := match.PlayerReadyDataMessage{}

	world.AddPlayerUnits(data1.UnitData, 0)
	world.AddPlayerUnits(data2.UnitData, 1)

	return world
}

func BenchmarkMatch(b *testing.B) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":11,\"y\":1}, {\"x\":11,\"y\":2}, {\"x\":11,\"y\":1}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":14,\"y\":1}]}}")

	world := CreateTestWorld(u1, u2)
	b.ResetTimer()
	for x := 0; x < 1000; x++ {
		world.Update()
	}
}

func BenchmarkGridUpdate(b *testing.B) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":11,\"y\":1}, {\"x\":11,\"y\":2}, {\"x\":11,\"y\":1}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":14,\"y\":1}]}}")

	world := CreateTestWorld(u1, u2)
	b.ResetTimer()
	for x := 0; x < 1000; x++ {
		world.Grid.Update()
	}
}
