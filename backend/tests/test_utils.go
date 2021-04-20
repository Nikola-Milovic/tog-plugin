package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/startup"
	"os"
	"strings"
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/match"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
)

var p1Units = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[{\"x\":6,\"y\":10}],\"knight\":[{\"x\":9,\"y\":7}, {\"x\":9,\"y\":3}]}}")
var p2Units = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":9,\"y\":10}]}}")

func CreateTestWorld(unitData []byte, unitData2 []byte, testing *testing.T) *game.World {
	world := game.CreateWorld()
	registry.RegisterWorld(world)

	world.AddPlayer("")
	world.AddPlayer("")

	data1 := match.PlayerReadyDataMessage{}
	if err := json.Unmarshal(unitData, &data1); err != nil {
		fmt.Printf("Error unmarshaling unitData %s", err.Error())
		testing.FailNow()
	}

	data2 := match.PlayerReadyDataMessage{}
	if err := json.Unmarshal(unitData2, &data2); err != nil {
		fmt.Printf("Error unmarshaling unitData %s", err.Error())
		testing.FailNow()
	}

	world.AddPlayerUnits(data1.UnitData, 0)
	world.AddPlayerUnits(data2.UnitData, 1)

	return world
}

func printImapsToFile() {
	f, err := os.Create("./temp.txt")
	check(err)

	defer f.Close()

	var sb strings.Builder
	for key, template := range startup.ProximityTemplates {
		heading := fmt.Sprintf("Proximity Map Key : %d, \n", key)
		sb.WriteString(heading)
		imap := template.Imap
		for y := 0; y < imap.Height; y++ {
			for x := 0; x < imap.Width; x++ {
				s := fmt.Sprintf("%.2f ", imap.Grid[x][y])
				sb.WriteString(s)
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	for key, template := range startup.InterestTemplates {
		heading := fmt.Sprintf("Interest Map Key : %d, \n", key)
		sb.WriteString(heading)
		imap := template.Imap
		for y := 0; y < imap.Height; y++ {
			for x := 0; x < imap.Width; x++ {
				s := fmt.Sprintf("%.2f ", imap.Grid[x][y])
				sb.WriteString(s)
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	w := bufio.NewWriter(f)
	writtenBytes, err := w.WriteString(sb.String())
	check(err)
	fmt.Printf("wrote %d bytes\n", writtenBytes)

	w.Flush()
}

func createTempFile() {
	_, err := os.Create("./temp.txt")
	check(err)
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
