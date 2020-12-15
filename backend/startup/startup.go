package startup

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Nikola-Milovic/tog-plugin/constants"
)

//StartUp is called when the server is started, all of the thing that should be done once server is started should be placed here
func StartUp() {
	println("StartUp")
	PopulateUnitDataMap()
}

//UnitData represents a map, where the key is the unitID and the value is the map[string]interface{} representing its data, components...
var UnitData = make(map[string]interface{}, 10)

//PopulateUnitDataMap populates the map with JSON data from the resources folder, it is executed once on server startup
//and is available for the rest of the server lifespan
func PopulateUnitDataMap() {

	path := "/nakama/data/units.json"

	jsonData, _ := ioutil.ReadFile(path)
	var data []interface{}

	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		panic("Couldn't unmarshal populateIndex")
	}

	for _, d := range data {
		unitData := d.(map[string]interface{})
		id := unitData[constants.UnitIDJson].(string)
		UnitData[id] = d
	}
}
