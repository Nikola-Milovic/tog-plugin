package ecs

import "github.com/Nikola-Milovic/tog-plugin/src/constants"

//Entity currently useless, just an index
type Entity struct {
	PlayerTag byte // 0 or 1, player 1 or 2
	Index     int
}

//EntityData represents the data the client receives
type EntityData struct {
	Position constants.V2 `json:"position"`
	State    string       `json:"state"`
	Index    int          `json:"index"`
}
