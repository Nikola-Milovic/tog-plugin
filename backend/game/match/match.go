package match

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/registry"
	"github.com/heroiclabs/nakama-common/runtime"
)

// OpCode represents a enum for valid OpCodes
// used by the match logic
type OpCode int64

const TICK_RATE = 5

const (
	//OpCodeClientEvents is used to indicate that this data is sending the current state of the game to the clients
	OpCodeClientEvents     = 1
	OpCodeMatchPreperation = 2
	OpCodeMatchEnd         = 3
	OpCodeMatchStart       = 4
	OpCodePlayerReady      = 5
)

// Match is the object registered
// as a runtime.Match interface
type Match struct{}

//Player is used by the match to keep track of the current state of the player and his actions
type Player struct {
	DisplayName string
	Ready       bool
	ID          string
	Tag         int
	PlayerWon   bool
}

// MatchData holds information that is passed between
// Nakama match methods
type MatchData struct {
	presences  map[string]runtime.Presence
	matchState MatchState
	World      *game.World
	Players    map[string]*Player
}

// GetPresenceList returns an array of current precenes in an array
func (state *MatchData) GetPresenceList() []runtime.Presence {
	presenceList := []runtime.Presence{}
	for _, presence := range state.presences {
		presenceList = append(presenceList, presence)
	}
	return presenceList
}

// MatchInit is called when a new match is created
func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	w := game.CreateWorld()

	matchData := &MatchData{
		presences:  map[string]runtime.Presence{},
		World:      w,
		matchState: MatchWaitingForPlayerState,
		Players:    make(map[string]*Player, 2),
	}
	tickRate := TICK_RATE
	label := "{\"name\": \"Game World\"}"

	registry.RegisterWorld(w)

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

	for _, presence := range presences {
		// Add the player that joined to the presence
		matchData.presences[presence.GetUserId()] = presence
		tag := matchData.World.AddPlayer(presence.GetUserId())
		matchData.Players[presence.GetUserId()] = &Player{Tag: tag, Ready: false, ID: presence.GetUserId(), DisplayName: presence.GetUsername()}
		logger.Info("Match joined %v\n", tag)
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
		matchEnd(matchData, logger, dispatcher, ctx, nk)
		fmt.Printf("Match end\n")
		return nil
	}

	//If we are still waiting for players to join the match, just return
	if matchData.matchState == MatchWaitingForPlayerState {
		logger.Info("Waiting for players state")

		//If there are 2 players, start the preperation state, and give each player their tag
		if len(matchData.GetPresenceList()) == 2 {
			matchPreperation(data, logger, dispatcher, ctx, nk)
		}

		return matchData
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

				//Unmarshal name
				playerReadyData := PlayerReadyDataMessage{}
				if err := json.Unmarshal(message.GetData(), &playerReadyData); err != nil {
					logger.Error("Couldnt unmarshal playerReadyData %v", err.Error())
				}

				matchData.Players[message.GetUserId()].DisplayName = playerReadyData.Name
				println(playerReadyData.Name)
				matchData.World.AddPlayerUnits(playerReadyData.UnitData, matchData.Players[message.GetUserId()].Tag)
				if checkIfAllPlayersReady(data) {
					matchStarted(data, logger, dispatcher)
				}
			}
		}
		return matchData
	}

	matchData.World.Update()

	//Get the events needed to recreate the state on clients ----------
	clientEvents, err := matchData.World.GetClientEvents()

	if err != nil {
		logger.Error("Error getting entities data %e", err.Error())
	} else if len(clientEvents) > 0 {
		if sendErr := dispatcher.BroadcastMessage(OpCodeClientEvents, clientEvents, matchData.GetPresenceList(), nil, true); sendErr != nil {
			logger.Error(sendErr.Error())
		}
	}
	// ------------------------------------------------
	// for _, message := range messages {

	// }

	if !matchData.World.MatchActive {
		m.checkMatchEnd(matchData, logger, dispatcher)
	}

	return matchData
}

// MatchTerminate is code that is executed when the match ends
func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, data interface{}, graceSeconds int) interface{} {
	return data
}
