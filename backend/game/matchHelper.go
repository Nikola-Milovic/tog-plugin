package game

import (
	"encoding/json"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
)

type PlayerJoinedResponse struct {
	Tag int `json:"tag"`
}

func playedJoinedResponse(tag int, presence runtime.Presence, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	playerJoinedMessage := PlayerJoinedResponse{Tag: tag}
	jsonData, err := json.Marshal(playerJoinedMessage)
	if err != nil {
		logger.Error(err.Error())
	}
	if sendErr := dispatcher.BroadcastMessage(OpCodePlayerJoined, jsonData, []runtime.Presence{presence}, nil, true); sendErr != nil {
		logger.Error(sendErr.Error())
	}
}

func checkIfAllPlayersReady(data interface{}) bool {
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo panic or terminate match
		fmt.Sprintln("Invalid match state on checkIfAllPlayersReady!")
		return false
	}

	trueCount := 0

	for _, presence := range matchData.GetPrecenseList() {
		if matchData.Players[presence.GetUserId()].Ready == true {
			trueCount++
		}
	}

	if trueCount == 2 {
		return true
	}

	return false
}

func (m *Match) checkForGameEnd(data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	//changeMatchState(MatchStartedState, data, logger, dispatcher)
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on matchStarted!")
	}

	if matchData.matchState == MatchEndState {
		return
	}

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

	if player0Lost || Player1Lost {
		matchData.World.MatchActive = false
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
