package match

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/math"

	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/heroiclabs/nakama-common/runtime"
)

type MatchState int

const (
	MatchStartedState          MatchState = 0
	MatchPreperationState      MatchState = 1
	MatchWaitingForPlayerState MatchState = 2
	MatchEndState              MatchState = 3
)

type PreperationStateMessage struct {
	Tag   int            `json:"tag"`
	Name  string         `json:"name"`
	Units map[string]int `json:"units"`
}

func matchPreperation(data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher, ctx context.Context, nk runtime.NakamaModule) {
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on matchPreperation!")
	}

	matchData.matchState = MatchPreperationState
	for _, presence := range matchData.GetPresenceList() {
		dataToSend := preperationStateData(matchData.Players[presence.GetUserId()].Tag, data, presence, logger, dispatcher, ctx, nk)
		if sendErr := dispatcher.BroadcastMessage(OpCodeMatchPreperation, dataToSend, []runtime.Presence{presence}, nil, true); sendErr != nil {
			logger.Error(sendErr.Error())
		}
		fmt.Printf("Preperation state for %v\n", string(dataToSend))

	}

	fmt.Println("PreparationState")
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

func matchEnd(data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher, ctx context.Context, nk runtime.NakamaModule) {
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on matchPreperation!")
	}

	p0Won := false
	p1Won := false
	for _, pl := range matchData.Players {
		if pl.PlayerWon {
			switch pl.Tag {
			case 0:
				p0Won = true
			case 1:
				p1Won = true
			}
		}
	}

	resultCode := MatchEndPlayer0Victory
	if p0Won && p1Won {
		resultCode = MatchEndDraw
	} else if p0Won {
		resultCode = MatchEndPlayer0Victory
	} else {
		resultCode = MatchEndPlayer1Victory
	}

	matchData.matchState = MatchEndState
	result := MatchEndResultMessage{Result: resultCode}
	jsonData, err := json.Marshal(result)
	if err != nil {
		logger.Error(err.Error())
	}

	matchData.matchState = MatchEndState
	if sendErr := dispatcher.BroadcastMessage(OpCodeMatchEnd, jsonData, matchData.GetPresenceList(), nil, true); sendErr != nil {
		logger.Error(sendErr.Error())
	}

	updateUserRunsOnMatchEnd(matchData, logger, ctx, nk)
}

//-------------------------------------------- MATCH START --------------------------------------------------------------

type MatchStartUnitDataMessage struct {
	Tag      int         `json:"tag"`
	UnitID   string      `json:"unit_id"`
	ID       int      `json:"id"`
	Index    int         `json:"index"`
	Position math.Vector `json:"position"`
}

func matchStarted(data interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
	matchData, ok := data.(*MatchData)
	if !ok {
		//Todo add somekind of error
		logger.Error("Invalid data on matchStarted!")
	}

	unitDataMessage := make([]MatchStartUnitDataMessage, 0, matchData.World.Players[0].NumberOfUnits+matchData.World.Players[1].NumberOfUnits)
	matchData.matchState = MatchStartedState
	for _, ent := range matchData.World.EntityManager.GetEntities() {
		unitData := MatchStartUnitDataMessage{}
		unitData.Index = ent.Index
		unitData.Tag = ent.PlayerTag
		unitData.UnitID = ent.UnitID
		unitData.ID = ent.ID
		unitData.Position = matchData.World.ObjectPool.Components["PositionComponent"][ent.Index].(components.PositionComponent).Position
		unitDataMessage = append(unitDataMessage, unitData)
	}

	messageJSON, err := json.Marshal(&unitDataMessage)

	if err != nil {
		panic(fmt.Sprintf("Cannot marshal unit data prior to match start %v\n", err.Error()))
	}

	if sendErr := dispatcher.BroadcastMessage(OpCodeMatchStart, messageJSON, matchData.GetPresenceList(), nil, true); sendErr != nil {
		logger.Error(sendErr.Error())
	}

	matchData.World.StartMatch()

	fmt.Printf("Match start\n")
}
