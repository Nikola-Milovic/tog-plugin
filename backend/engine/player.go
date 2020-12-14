package engine

type PlayerData struct {
	Tag           int
	NumberOfUnits int
}

type Player struct {
	Ready bool
	ID    string
	Tag   int
}
