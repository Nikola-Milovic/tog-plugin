package ecs

import "github.com/Nikola-Milovic/tog-plugin/src/constants"

//Entity currently useless, just an index
type Entity struct {
	Index int
}

//EntityData represents the data the client receives
type EntityData struct {
	Position constants.V2 `json:"position"`
	Action   string       `json:"action"`
	Index    int          `json:"index"`
}
