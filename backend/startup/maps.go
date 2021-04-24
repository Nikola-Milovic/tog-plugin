package startup

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

var ProximityTemplates = make([]*engine.ImapTemplate, 0, 6)
var InterestTemplates = make([]*engine.ImapTemplate, 0, 9)
var SizeTemplates = make(map[string]*engine.ImapTemplate, 5)

func linearCalc(dist, maxDist int, value float32) float32 {
	return value - (value * (float32(dist) / float32(maxDist+1)))
}

func LinearCalcHigherFalloff(dist, maxDist int, value float32) float32 {
	return value - (value * (float32(dist) / (float32(maxDist+1) * 1.2)))
}

func NoFallOffCalc(dist, maxDist int, value float32) float32 {
	return value
}

func EndFallOffCalc(dist, maxDist int, value float32) float32 {
	if dist == maxDist {
		return value / 2
	}
	return value
}

func initProxImapTemplates() {
	for i := 0; i < 20; i += 4 {
		size := (constants.TickRate*5*i)/constants.TileSize + 9
		imap := engine.NewImap(size, size, constants.TileSize)
		imap.PropagateInfluence(size/2, size/2, size, LinearCalcHigherFalloff, 1)
		template := engine.ImapTemplate{Radius: i, Type: constants.ImapTypeProximity, Imap: imap}
		ProximityTemplates = append(ProximityTemplates, &template)
	}
}

func initInterestImapsTemplates() { // TODO change the 80 and stuff
	for i := 4; i < 20; i += 4 {
		size := 10*i/constants.TileSize + 9
		imap := engine.NewImap(size, size, constants.TileSize)
		imap.PropagateInfluence(size/2, size/2, size, linearCalc, 1)
		template := engine.ImapTemplate{Radius: i, Type: constants.ImapTypeProximity, Imap: imap}
		InterestTemplates = append(InterestTemplates, &template)
	}
}

func initSizeImapsTemplates() {
	size32 := 32/constants.TileSize + 4
	imap32 := engine.NewImap(size32, size32, constants.TileSize)
	imap32.PropagateInfluence(size32/2, size32/2, size32, EndFallOffCalc, 1)
	template32 := engine.ImapTemplate{Radius: size32, Type: constants.ImapTypeProximity, Imap: imap32}
	SizeTemplates[constants.StandardSize] = &template32
}
