package engine

//Entity currently useless, just an index
type Entity struct {
	//	PlayerTag byte // 0 or 1, player 1 or 2
	Index int
	Name  string
	//	Size      Vector
	//	State     string
}

//EntityData represents the data the client receives
type EntityData struct {
	Position Vector   `json:"position"`
	State    string   `json:"state"`
	Index    int      `json:"index"`
	Path     []Vector `json:"path"`
	Tag      byte     `json:"player_tag"`
}
