package match

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/heroiclabs/nakama-common/runtime"
)

// OpCode represents a enum for valid OpCodes
// used by the match logic
type OpCode int64

const TICK_RATE = 5

const (
	//OpCodeUpdateEntities is used to indicate that this data is sending the current state of the game to the clients
	OpCodeUpdateEntities   = 1
	OpCodePlayerJoined     = 2
	OpCodeMatchPreperation = 3
	OpCodeMatchEnd         = 4
	OpCodeMatchStart       = 5
	OpCodePlayerReady      = 6
)

// Match is the object registered
// as a runtime.Match interface
type Match struct{}

//Player is used by the match to keep track of the current state of the player and his actions
type Player struct {
	Ready bool
	ID    string
	Tag   int
}

// MatchData holds information that is passed between
// Nakama match methods
type MatchData struct {
	presences  map[string]runtime.Presence
	matchState MatchState
	World      *game.World
	Players    map[string]*Player
}

// GetPrecenseList returns an array of current precenes in an array
func (state *MatchData) GetPrecenseList() []runtime.Presence {
	precenseList := []runtime.Presence{}
	for _, precense := range state.presences {
		precenseList = append(precenseList, precense)
	}
	return precenseList
}

// MatchInit is called when a new match is created
func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	matchData := &MatchData{
		presences:  map[string]runtime.Presence{},
		World:      game.CreateWorld(),
		matchState: MatchWaitingForPlayerState,
		Players:    make(map[string]*Player, 2),
	}
	tickRate := TICK_RATE
	label := "{\"name\": \"Game World\"}"

	logger.Info("Match created")

	return matchData, tickRate, label
}

// MatchJoinAttempt is called when a player tried to join a match
// and validates their attempt
func (m *Match) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	matchData, ok := state.(*MatchData)
	if !ok {
		logger.Error("Invalid match state on join attempt!")
		return state, false, "Invalid match state!"
	}

	// Validate user is not already connected
	if _, ok := matchData.presences[presence.GetUserId()]; ok {
		return matchData, false, "User already logged in."
	}
	return matchData, true, ""

}

// MatchJoin is called when a player successfully joins the match
func (m *Match) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, data interface{}, presences []runtime.Presence) interface{} {
	matchData, ok := data.(*MatchData)
	if !ok {
		logger.Error("Invalid match state on join!")
		return data
	}

	logger.Info("Match joined")

	for _, precense := range presences {
		// Add the player that joined to the presences
		matchData.presences[precense.GetUserId()] = precense
	}

	return matchData
}

// MatchLeave is called when a player leaves the match
func (m *Match) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, data interface{}, presences []runtime.Presence) interface{} {
	matchData, ok := data.(*MatchData)
	if !ok {
		logger.Error("Invalid match state on leave!")
		return data
	}
	for _, presence := range presences {
		if _, ok := matchData.presences[presence.GetUserId()]; ok {
			delete(matchData.presences, presence.GetUserId())
		}
		if _, ok := matchData.Players[presence.GetUserId()]; ok {
			delete(matchData.Players, presence.GetUserId())
		}
	}

	if len(matchData.presences) == 0 {
		logger.Info("Match Terminated")
		return nil
	}

	return matchData
}

// MatchLoop is code that is executed every tick
func (m *Match) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, data interface{}, messages []runtime.MatchData) interface{} {
	matchData, ok := data.(*MatchData)
	if !ok {
		logger.Error("Invalid match state on match loop!")
		return data
	}

	//If the match ended, terminate the match, should do some calculations and such here, mmr gain and rewards
	if matchData.matchState == MatchEndState {
		fmt.Printf("Match end\n")
		return nil
	}

	//If we are still waiting for players to join the match, just return
	if matchData.matchState == MatchWaitingForPlayerState {
		//	logger.Info("Waiting for players state")

		//If there are 2 players, start the preperation state, and give each player their tag
		if len(matchData.GetPrecenseList()) == 2 {
			for _, precense := range matchData.GetPrecenseList() {
				tag := matchData.World.AddPlayer()
				matchData.Players[precense.GetUserId()] = &Player{Tag: tag, Ready: false, ID: precense.GetUserId()}
				playedJoinedResponse(tag, precense, logger, dispatcher)
			}
			matchPreperation(data, logger, dispatcher)
		}

		return data
	}

	//Wait until both players confirm their armies
	if matchData.matchState == MatchPreperationState {
		for _, message := range messages {
			switch message.GetOpCode() {
			case OpCodePlayerReady:
				if matchData.Players[message.GetUserId()].Ready {
					fmt.Printf("Player already is ready\n")
					return matchData
				}
				fmt.Println("PlayerReady")
				matchData.Players[message.GetUserId()].Ready = true
				matchData.World.AddPlayerUnits(message.GetData(), matchData.Players[message.GetUserId()].Tag)
				if checkIfAllPlayersReady(data) {
					matchStarted(data, logger, dispatcher)
				}
			}
		}
		return matchData
	}

	matchData.World.Update()

	entityData, err := game.GetEntitiesData(matchData.World)

	if err != nil {
		logger.Error("Error getting entities data %e", err.Error())
	} else {
		if sendErr := dispatcher.BroadcastMessage(OpCodeUpdateEntities, entityData, matchData.GetPrecenseList(), nil, true); sendErr != nil {
			logger.Error(sendErr.Error())
		}
	}

	m.checkForGameEnd(matchData, logger, dispatcher)

	// for _, message := range messages {

	// }

	return matchData
}

// MatchTerminate is code that is executed when the match ends
func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, data interface{}, graceSeconds int) interface{} {
	return data
}
