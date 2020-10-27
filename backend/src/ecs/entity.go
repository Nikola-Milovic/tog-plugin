package ecs

import "github.com/Nikola-Milovic/tog-plugin/src/constants"

type Entity struct {
	Index int
}

type EntityData struct {
	Position constants.V2 `json:"position"`
	Action   string       `json:"action"`
	Index    int          `json:"index"`
}
