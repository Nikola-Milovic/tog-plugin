package tests

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/startup"
	"testing"
)

func TestWriting(t *testing.T) {
	printImapsToFile()
}

func TestSmallerIntoBigger(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":8}]}}")

	world := CreateTestWorld(u1, u2, t)
	grid := world.Grid

	imap := grid.GetEnemyProximityImap(0)

	newImap := engine.AddIntoBiggerMap(startup.ProximityTemplates[constants.MovementSpeedMedium].Imap, imap,
		10, 10, 1)

	printImapToFile(newImap, "AddIntoBiggerMap", false)
}

func TestBiggerIntoSmaller(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":8}]}}")

	CreateTestWorld(u1, u2, t)
//	grid := world.Grid

	bigger := startup.ProximityTemplates[constants.MovementSpeedFast].Imap
	smaller := startup.ProximityTemplates[constants.MovementSpeedSlow].Imap
	printImapToFile(smaller, "Before", false)
	printImapToFile(bigger, "Bigger", true)
	newImap := engine.AddIntoSmallerMap(bigger, smaller, 6, 6, -1)

	printImapToFile(newImap, "BiggerIntoSmallerMap", true)
}
