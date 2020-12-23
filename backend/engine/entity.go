package engine

//Entity just a holder, represents an index, hold the player tag, its ID and whether its active or not
type Entity struct {
	PlayerTag int // 0 or 1, player 1 or 2
	Index     int
	ID        string
	Active    bool
}

//NewEntityData represents a struct that holds data needed to add a new entity
type NewEntityData struct {
	PlayerTag int
	ID        string
	Data      interface{}
}

//EntityMessage represents the data the client receives
type EntityMessage struct {
	Position Vector `json:"position"`
	State    string `json:"state"`
	Index    int    `json:"index"`
	//	Path     []Vector `json:"path"`
	Tag    int `json:"player_tag"`
	Health int `json:"health"`
}
