package engine

//Entity just a holder, represents an index, hold the player tag, its ID and whether its active or not
// the ID is used for caching, so we can quickly check if the target is still valid
type Entity struct {
	PlayerTag int // 0 or 1, player 1 or 2
	Index     int
	UnitID    string
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
	ID       string `json:"id"`
	Health int `json:"health"`
}
