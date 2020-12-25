package startup

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Nikola-Milovic/tog-plugin/constants"
)

//StartUp is called when the server is started, all of the thing that should be done once server is started should be placed here
func StartUp(testing bool) {
	println("StartUp")
	PopulateUnitDataMap(testing)
	PopulateEffectDataMap(testing)
	PopulateAbilityDataMap(testing)
}

//UnitDataMap represents a map, where the key is the unitID and the value is the map[string]interface{} representing its data, components...
var UnitDataMap = make(map[string]map[string]interface{}, 10)

//PopulateUnitDataMap populates the map with JSON data from the resources folder, it is executed once on server startup
//and is available for the rest of the server lifespan
func PopulateUnitDataMap(testing bool) {

	path := "/nakama/resources/units.json"

	if testing {
		path = "../resources/units.json"
	}

	jsonData, _ := ioutil.ReadFile(path)
	var data []interface{}

	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		panic("Couldn't unmarshal UnitDataMap")
	}

	for _, d := range data {
		unitData := d.(map[string]interface{})
		id := unitData[constants.UnitIDJson].(string)
		UnitDataMap[id] = unitData
	}
}

//EffectDataMap represents
var EffectDataMap = make(map[string]map[string]interface{}, 10)

//PopulateEffectDataMap populates
func PopulateEffectDataMap(testing bool) {

	path := "/nakama/resources/effects.json"

	if testing {
		path = "../resources/effects.json"
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
		EffectDataMap[id] = effData
	}
}

//AbilityDataMap represents
var AbilityDataMap = make(map[string]map[string]interface{}, 10)

//PopulateAbilityDataMap is
func PopulateAbilityDataMap(testing bool) {

	path := "/nakama/resources/abilities.json"

	if testing {
		path = "../resources/abilities.json"
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
		AbilityDataMap[id] = abData
	}
}
