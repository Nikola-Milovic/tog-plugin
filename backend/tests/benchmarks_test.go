package tests

import (
	"encoding/json"
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/startup"
	"io/ioutil"
	"testing"
	"time"
	"unsafe"
)

func TestEntityIDMemory(t *testing.T) {
	s := "-NDveu-9Q"
	fmt.Println("Size of id:", unsafe.Sizeof(s))
}

func TestSizeOfGlobalVariables(t *testing.T) {
	fmt.Printf("Size of unit data map is %v\n", Of(constants.UnitDataMap))
	fmt.Printf("Size of ability data map is %v\n", Of(constants.AbilityDataMap))
	fmt.Printf("Size of effect data map is %v\n", Of(constants.EffectDataMap))
	fmt.Printf("Size of proximity templates is %v\n", Of(startup.ProximityTemplates))
	fmt.Printf("Size of interest templates is %v\n", Of(startup.InterestTemplates))
}


func TestWorldSize_WithTwoEntities(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := CreateTestWorld(p1Units, p2Units, t)

	fmt.Println("WITH THE 2 ENTITIES")

	fmt.Printf("Size of EntityManager.Entities is %v\n", Of(world.EntityManager.GetEntities()))
	fmt.Printf("Size of All components is %v\n", Of(world.ObjectPool.Components))
	fmt.Printf("Size of EventManager is %v\n", Of(world.EventManager))
	fmt.Printf("Size of Grid is %v\n", Of(world.Grid))
	fmt.Printf("Size of players is %v\n", Of(world.Players))
	fmt.Printf("Size of AI's is %v\n", Of(world.ObjectPool.AI))

	totalSize := Of(world.EntityManager.GetEntities()) + Of(world.ObjectPool.Components) + Of(world.EventManager) +
		Of(world.Grid) + Of(world.Players) + Of(world.ObjectPool.AI)


	fmt.Printf("\n\n TOTAL SIZE IN BYTES IS : %v\n\n", totalSize)
}

func TestWorldSize_WithEightEntities(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := CreateTestWorld(p1Units, p2Units, t)

	fmt.Println("WITH THE 8 ENTITIES")

	fmt.Printf("Size of EntityManager.Entities is %v\n", Of(world.EntityManager.GetEntities()))
	fmt.Printf("Size of All components is %v\n", Of(world.ObjectPool.Components))
	fmt.Printf("Size of EventManager is %v\n", Of(world.EventManager))
	fmt.Printf("Size of Grid is %v\n", Of(world.Grid))
	fmt.Printf("Size of players is %v\n", Of(world.Players))
	fmt.Printf("Size of AI's is %v\n", Of(world.ObjectPool.AI))

	totalSize := Of(world.EntityManager.GetEntities()) + Of(world.ObjectPool.Components) + Of(world.EventManager) +
		Of(world.Grid) + Of(world.Players) + Of(world.ObjectPool.AI)

	fmt.Printf("\n\n TOTAL SIZE IN BYTES IS : %v\n\n", totalSize)
}

func TestTickSpeed_CreateWorld(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}
	start := time.Now()

	CreateTestWorld(p1Units, p2Units, t)

	//world.Update()

	end := time.Now()
	fmt.Printf("\n\n Total time taken to create world : %v\n\n", end.Sub(start))
}

func TestTickSpeed_SingleTick(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}

	world := CreateTestWorld(p1Units, p2Units, t)

	start := time.Now()

	world.Update()

	end := time.Now()
	fmt.Printf("\n\n Total time taken for single world update is : %v\n\n", end.Sub(start))
}

func TestSpeed_MatchWith8Entities(t *testing.T) {
	jsonData, _ := ioutil.ReadFile("../resources/units.json")
	var data []map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Errorf("Couldn't unmarshal json: %e", err)
	}
	world := CreateTestWorld(p1Units, p2Units, t)

	start := time.Now()

	for world.MatchActive {
		world.Update()
	}

	end := time.Now()
	fmt.Printf("\n\n Total time taken for 100 world updates is : %v\n\n", end.Sub(start))
}
