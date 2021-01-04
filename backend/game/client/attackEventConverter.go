package client

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

func AttackEventConverter(ev engine.Event) map[string]interface{} {
	data := make(map[string]interface{}, 3)

	data["event"] = "attack"
	data["me"] = ev.Data["emitter"]
	data["who"] = ev.Data["target"]

	return data
}

// data := make(map[string]interface{}, 1)
// data["target"] = indx
// ev := engine.Event{Index: index, ID: "AttackEvent", Priority: 100, Data: data}
