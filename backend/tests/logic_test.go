package tests

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/math"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
)

func TestEightEntitiesFighting(t *testing.T) {
	world := CreateTestWorld(p1Units, p2Units, t)

	p1 := world.ObjectPool.Components["PositionComponent"][0].(components.PositionComponent)
	p2 := world.ObjectPool.Components["PositionComponent"][1].(components.PositionComponent)
	p3 := world.ObjectPool.Components["PositionComponent"][2].(components.PositionComponent)
	p4 := world.ObjectPool.Components["PositionComponent"][3].(components.PositionComponent)
	p5 := world.ObjectPool.Components["PositionComponent"][4].(components.PositionComponent)
	p6 := world.ObjectPool.Components["PositionComponent"][5].(components.PositionComponent)
	p7 := world.ObjectPool.Components["PositionComponent"][6].(components.PositionComponent)
	p8 := world.ObjectPool.Components["PositionComponent"][7].(components.PositionComponent)

	p1.Position = math.Vector{0, 5}
	p2.Position = math.Vector{10, 30}
	p5.Position = math.Vector{60, 80}
	p6.Position = math.Vector{90, 40}
	p4.Position = math.Vector{190, 35}
	p7.Position = math.Vector{60, 50}
	p8.Position = math.Vector{10, 30}
	p3.Position = math.Vector{150, 40}

	world.ObjectPool.Components["PositionComponent"][0] = p1
	world.ObjectPool.Components["PositionComponent"][1] = p2
	world.ObjectPool.Components["PositionComponent"][2] = p3
	world.ObjectPool.Components["PositionComponent"][3] = p4
	world.ObjectPool.Components["PositionComponent"][4] = p5
	world.ObjectPool.Components["PositionComponent"][5] = p6
	world.ObjectPool.Components["PositionComponent"][6] = p7
	world.ObjectPool.Components["PositionComponent"][7] = p8

	h1 := world.ObjectPool.Components["StatsComponent"][1].(components.StatsComponent)
	h1.Health = 30

	h5 := world.ObjectPool.Components["StatsComponent"][5].(components.StatsComponent)
	h5.Health = 18

	h7 := world.ObjectPool.Components["StatsComponent"][7].(components.StatsComponent)
	h7.Health = 45

	world.ObjectPool.Components["StatsComponent"][1] = h1
	world.ObjectPool.Components["StatsComponent"][5] = h5
	world.ObjectPool.Components["StatsComponent"][7] = h7

	for world.MatchActive {
		world.Update()
	}
}

func TestTwoEntities(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5},{\"x\":5,\"y\":4}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":8}]}}")

	world := CreateTestWorld(u1, u2, t)

	for x := 0; x < 200; x++ {
		world.Update()
		if x%5 == 0 {
			engine.PrintImapToFileWithStep(world.Grid.GetOccupationalMap(), fmt.Sprintf("Second is %d", x/5), 4)
		}
	}
}

func TestFourEntities(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":4,\"y\":5}, {\"x\":5,\"y\":5}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5},  {\"x\":4,\"y\":5},  {\"x\":3,\"y\":6}]}}")

	world := CreateTestWorld(u1, u2, t)

	for x := 0; x < 200; x++ {
		world.Update()
		if x%5 == 0 {
			engine.PrintImapToFileWithStep(world.Grid.GetOccupationalMap(), fmt.Sprintf("Second is %d", x/5), 1)
		}
	}
}

func TestEightEntities(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":4,\"y\":5}, {\"x\":5,\"y\":5}, {\"x\":5,\"y\":6}, {\"x\":4,\"y\":5}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":4,\"y\":4}, {\"x\":5,\"y\":5}, {\"x\":5,\"y\":6}, {\"x\":4,\"y\":5}]}}")

	world := CreateTestWorld(u1, u2, t)

	for x := 0; x < 200; x++ {
		world.Update()
		if x%5 == 0 {
			engine.PrintImapToFileWithStep(world.Grid.GetOccupationalMap(), fmt.Sprintf("Second is %d", x/5), 1)
		}
	}
}
