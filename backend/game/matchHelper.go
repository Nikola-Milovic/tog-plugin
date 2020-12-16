package game

import (
	"encoding/json"
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/heroiclabs/nakama-common/runtime"
)

type PlayerJoinedResponse struct {
	Tag int `json:"tag"`
}

func playedJoinedResponse(tag int, presence runtime.Presence, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	playerJoinedMessage := PlayerJoinedResponse{Tag: 0}
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

func changeMatchState(newState MatchState, data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on changeMatchState!")
	}
	matchData.matchState = newState
	matchStateData := MatchStateMessage{MatchState: newState}
	jsonData, err := json.Marshal(matchStateData)
	if err != nil {
		logger.Error(err.Error())
	}
	if sendErr := dispatcher.BroadcastMessage(OpCodeMatchStateChange, jsonData, matchData.GetPrecenseList(), nil, true); sendErr != nil {
		logger.Error(sendErr.Error())
	}
	fmt.Printf("New state of the match is %v", newState)
}

type MatchStartUnitDataMessage struct {
	Tag      int           `json:"tag"`
	UnitID   string        `json:"id"`
	Index    int           `json:"index"`
	Position engine.Vector `json:"position"`
}

func (m *Match) matchStarted(data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	changeMatchState(MatchStartedState, data, logger, dispatcher)
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on matchStarted!")
	}

	unitDataMessage := make([]MatchStartUnitDataMessage, 0, matchData.World.Players[0].NumberOfUnits+matchData.World.Players[1].NumberOfUnits)

	for _, ent := range matchData.World.EntityManager.Entities {
		if ent.Active {
			unitData := MatchStartUnitDataMessage{}
			unitData.Index = ent.Index
			unitData.Tag = ent.PlayerTag
			unitData.UnitID = ent.ID
			unitData.Position = matchData.World.ObjectPool.Components["PositionComponent"][ent.Index].(PositionComponent).Position
			unitDataMessage = append(unitDataMessage, unitData)
		}
	}

	messageJSON, err := json.Marshal(&unitDataMessage)

	if err != nil {
		panic(fmt.Sprintf("Cannot marshal unit data prior to match start %v\n", err.Error()))
	}

	if sendErr := dispatcher.BroadcastMessage(OpCodeMatchStart, messageJSON, matchData.GetPrecenseList(), nil, true); sendErr != nil {
		logger.Error(sendErr.Error())
	}

}
