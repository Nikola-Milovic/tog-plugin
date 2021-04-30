package startup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	rand2 "math/rand"
	"time"

	"github.com/Nikola-Milovic/tog-plugin/constants"
)

//StartUp is called when the server is started, all of the thing that should be done once server is started should be placed here
func StartUp(testing bool) {
	println("StartUp")
	populateUnitDataMap(testing)
	populateEffectDataMap(testing)
	populateAbilityDataMap(testing)
	initProxImapTemplates()
	initInterestImapsTemplates()
	initSizeImapsTemplates()

	rand2.Seed(time.Now().UnixNano())
}

var ResourcesPath = "../resources"

//populateUnitDataMap populates the map with JSON data from the resources folder, it is executed once on server startup
//and is available for the rest of the server lifespan
func populateUnitDataMap(testing bool) {

	path := "/nakama/resources/units.json"

	if testing {
		path = fmt.Sprintf("%s/units.json", ResourcesPath)
	}

	jsonData, _ := ioutil.ReadFile(path)
	var data []interface{}

	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		panic(fmt.Sprintf("Couldn't unmarshal UnitDataMap, %v", err.Error()))
	}

	for _, d := range data {
		unitData := d.(map[string]interface{})
		id := unitData[constants.UnitIDJson].(string)
		constants.UnitDataMap[id] = unitData
	}
}

//populateEffectDataMap populates
func populateEffectDataMap(testing bool) {

	path := "/nakama/resources/effects.json"

	if testing {
		path = fmt.Sprintf("%s/effects.json", ResourcesPath)
	}

	jsonData, _ := ioutil.ReadFile(path)
	var data []interface{}

	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		panic("Couldn't unmarshal EffectDataMap")
	}

	for _, d := range data {
		effData := d.(map[string]interface{})
		id := effData[constants.EffectIDJson].(string)
		constants.EffectDataMap[id] = effData
	}
}

//populateAbilityDataMap is
func populateAbilityDataMap(testing bool) {

	path := "/nakama/resources/abilities.json"

	if testing {
		path = fmt.Sprintf("%s/abilities.json", ResourcesPath)
	}

	jsonData, _ := ioutil.ReadFile(path)
	var data []interface{}

	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		panic("Couldn't unmarshal AbilityDataMap")
	}

	for _, d := range data {
		abData := d.(map[string]interface{})
		id := abData[constants.AbilityIDJson].(string)
		constants.AbilityDataMap[id] = abData
	}
}
