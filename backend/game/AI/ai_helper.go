package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
)

func canActivateAbility(lastActivated int, abilityID string, w *game.World) bool {

	if w.Tick-lastActivated > int(constants.AbilityDataMap[abilityID]["Cooldown"].(float64)) {
		return true
	}

	return false
}
