package match

import (
	"encoding/json"
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/heroiclabs/nakama-common/runtime"
)

type MatchState int

const (
	MatchStartedState          MatchState = 0
	MatchPreperationState      MatchState = 1
	MatchWaitingForPlayerState MatchState = 2
	MatchEndState              MatchState = 3
)

func matchPreperation(data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on matchPreperation!")
	}
	matchData.matchState = MatchPreperationState
	if sendErr := dispatcher.BroadcastMessage(OpCodeMatchPreperation, nil, matchData.GetPrecenseList(), nil, true); sendErr != nil {
		logger.Error(sendErr.Error())
	}
	fmt.Printf("PreparationState")
}

//--------------------------------------------------- MATCH END ---------------------------------------------------------------

//MatchEndResultMessage is sent to clients when the match ends, the result is 1, 0 or 2.
//0 is player 1 victory, 1 is player 2 victory, 2 is a tie
type MatchEndResultMessage struct {
	Result int `json:"result"`
}

const (
	MatchEndPlayer0Victory = 0
	MatchEndPlayer1Victory = 1
	MatchEndDraw           = 2
)

func matchEnd(resultCode int, data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on matchPreperation!")
	}
	matchData.matchState = MatchEndState
	fmt.Println("Player 0 lost")
	result := MatchEndResultMessage{Result: resultCode}
	jsonData, err := json.Marshal(result)
	if err != nil {
		logger.Error(err.Error())
	}

	matchData.matchState = MatchEndState
	if sendErr := dispatcher.BroadcastMessage(OpCodeMatchEnd, jsonData, matchData.GetPrecenseList(), nil, true); sendErr != nil {
		logger.Error(sendErr.Error())
	}
}

//-------------------------------------------- MATCH START --------------------------------------------------------------

type MatchStartUnitDataMessage struct {
	Tag      int           `json:"tag"`
	UnitID   string        `json:"id"`
	Index    int           `json:"index"`
	Position engine.Vector `json:"position"`
}

func matchStarted(data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on matchStarted!")
	}

	unitDataMessage := make([]MatchStartUnitDataMessage, 0, matchData.World.Players[0].NumberOfUnits+matchData.World.Players[1].NumberOfUnits)
	matchData.matchState = MatchStartedState
	for _, ent := range matchData.World.EntityManager.Entities {
		if ent.Active {
			unitData := MatchStartUnitDataMessage{}
			unitData.Index = ent.Index
			unitData.Tag = ent.PlayerTag
			unitData.UnitID = ent.ID
			unitData.Position = matchData.World.ObjectPool.Components["PositionComponent"][ent.Index].(game.PositionComponent).Position
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

	fmt.Printf("Match start")
}
