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
	radius16 := 32/constants.TileSize + 5
	imap16 := engine.NewImap(radius16, radius16, constants.TileSize)
	imap16.PropagateInfluence(radius16/2, radius16/2, radius16, EndFallOffCalc, 1)
	template16 := engine.ImapTemplate{Radius: radius16, Type: constants.ImapTypeProximity, Imap: imap16}
	SizeTemplates[constants.StandardSize] = &template16

	radius10 := 10/constants.TileSize + 5
	imap10 := engine.NewImap(radius10, radius10, constants.TileSize)
	imap10.PropagateInfluence(radius10/2, radius10/2, radius10, EndFallOffCalc, 1)
	template10 := engine.ImapTemplate{Radius: radius10, Type: constants.ImapTypeProximity, Imap: imap10}
	SizeTemplates[constants.SmallSize] = &template10
}
