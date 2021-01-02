package match

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/heroiclabs/nakama-common/runtime"
)

type PlayerReadyDataMessage struct {
	Name     string                     `json:"name"`
	UnitData map[string][]engine.Vector `json:"units"`
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
					playerJoinedMessage.Units[unit] = len(runData.Army[unit])
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
					playerJoinedMessage.Units[unit] = len(runData.Army[unit])
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

func (m *Match) matchEnd(data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	//changeMatchState(MatchStartedState, data, logger, dispatcher)
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on matchStarted!")
	}

	matchData.matchState = MatchEndState

	player0Lost := false
	Player1Lost := false
	for _, player := range matchData.World.Players {
		if player.NumberOfUnits == 0 {
			switch player.Tag {
			case 0:
				player0Lost = true
			case 1:
				Player1Lost = true
			}
		}
	}

	if player0Lost && Player1Lost { // Draw/ Tie
		fmt.Println("Draw")
		matchEnd(MatchEndDraw, matchData, logger, dispatcher)
		return
	}

	if player0Lost { //Player with tag 0 lost
		fmt.Println("Player 1 victory")
		matchEnd(MatchEndPlayer1Victory, matchData, logger, dispatcher)
		return
	}

	if Player1Lost { // player with tag 1 lost
		fmt.Println("Player 0 victory")
		matchEnd(MatchEndPlayer0Victory, matchData, logger, dispatcher)
		return
	}
}
