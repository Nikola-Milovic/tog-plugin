package startup

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
)

var ProximityTemplates = make(map[float32]*engine.ImapTemplate, 3)
var InterestTemplates = make(map[int]*engine.ImapTemplate, 10)

func linearCalc(dist, maxDist int, value float32) float32 {
	return value - (value * (float32(dist) / float32(maxDist+1)))
}

func initProxMapTemplates() {
	speeds := []float32{constants.MovementSpeedSlow, constants.MovementSpeedMedium, constants.MovementSpeedFast}
	for i := 0; i < 3; i++ {
		size := 2 * (int(speeds[i])/constants.TileSize) + 1
		imap := engine.NewImap(size, size, constants.TileSize)
		imap.PropagateInfluence(size/2, size/2, size, linearCalc, 1)
		template := engine.ImapTemplate{Radius: i, Type: constants.ImapTypeProximity, Imap: imap}
		ProximityTemplates[speeds[i]] = &template
	}
}

func initInterestImapsTemplates() { // TODO change the 80 and stuff
	for i := 0; i < 320; i+= 80 {
		size := 2 * i/constants.TileSize + 1
		imap := engine.NewImap(size, size, constants.TileSize)
		imap.PropagateInfluence(size/2, size/2, size, linearCalc, 1)
		template := engine.ImapTemplate{Radius: i, Type: constants.ImapTypeProximity, Imap: imap}
		InterestTemplates[i] = &template
	}
}