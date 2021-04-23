package match

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/math"

	"github.com/heroiclabs/nakama-common/runtime"
)

type PlayerReadyDataMessage struct {
	Name     string                   `json:"name"`
	UnitData map[string][]math.Vector `json:"units"`
}

func preperationStateData(tag int, data interface{}, presence runtime.Presence, logger runtime.Logger, dispatcher runtime.MatchDispatcher, ctx context.Context, nk runtime.NakamaModule) []byte {
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo panic or terminate match
		fmt.Sprintln("Invalid match state on playedJoinedResponse!")
	}

	playerJoinedMessage := PreperationStateMessage{Tag: tag}
	switch tag {
	case 0:
		for _, pl := range matchData.Players {
			if pl.Tag == 1 {
				playerJoinedMessage.Name = pl.DisplayName

				runData, _ := getUserRun(pl.ID, ctx, logger, nk)
				playerJoinedMessage.Units = make(map[string]int)
				for unit := range runData.Army {
					playerJoinedMessage.Units[unit] = runData.Army[unit]
				}
			}
		}
	case 1:
		for _, pl := range matchData.Players {
			if pl.Tag == 0 {
				playerJoinedMessage.Name = pl.DisplayName

				runData, _ := getUserRun(pl.ID, ctx, logger, nk)
				playerJoinedMessage.Units = make(map[string]int)
				for unit := range runData.Army {
					playerJoinedMessage.Units[unit] = runData.Army[unit]
				}
			}
		}
	}
	jsonData, err := json.Marshal(playerJoinedMessage)
	if err != nil {
		logger.Error(err.Error())
	}

	return jsonData
}

func checkIfAllPlayersReady(data interface{}) bool {
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo panic or terminate match
		fmt.Sprintln("Invalid match state on checkIfAllPlayersReady!")
		return false
	}

	trueCount := 0

	for _, presence := range matchData.GetPresenceList() {
		if matchData.Players[presence.GetUserId()].Ready == true {
			trueCount++
		}
	}

	if trueCount == 2 {
		return true
	}

	return false
}

func (m *Match) checkMatchEnd(data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	//changeMatchState(MatchStartedState, data, logger, dispatcher)
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on matchStarted!")
	}

	matchData.matchState = MatchEndState
	for _, player := range matchData.World.Players {
		if player.NumberOfUnits == 0 {
			switch player.Tag {
			case 0:
				matchData.Players[matchData.World.Players[1].ID].PlayerWon = true
			case 1:
				matchData.Players[matchData.World.Players[0].ID].PlayerWon = true
			}
		}
	}
}

func updateUserRunsOnMatchEnd(mData interface{}, logger runtime.Logger, ctx context.Context, nk runtime.NakamaModule) {
	matchData, ok := mData.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on updateUserRunsOnMatchEnd!")
	}

	for _, player := range matchData.Players {

		playerWon := player.PlayerWon
		runData, _ := getUserRun(player.ID, ctx, logger, nk)

		if playerWon {
			runData.Floor++
			runData.Score.Wins++
		} else {
			runData.Score.Losses++
		}

		//Write the user run data back
		runDataJSON, err := json.Marshal(runData)

		if err != nil {
			logger.Error("Error marshaling run data %v", err.Error())
		}

		writeObjects := []*runtime.StorageWrite{
			{
				Collection:      "runs",
				Key:             "current_run",
				UserID:          player.ID,
				Value:           string(runDataJSON),
				PermissionRead:  1,
				PermissionWrite: 0,
						},
		}

		if _, err := nk.StorageWrite(ctx, writeObjects); err != nil {
			logger.Error("Error writing the run data after altering it in match end", err.Error())
		}

	}

}
