package control

import (
	"context"
	"database/sql"

	"github.com/Nikola-Milovic/tog-plugin/src/ecs"
	"github.com/heroiclabs/nakama-common/runtime"
)

// OpCode represents a enum for valid OpCodes
// used by the match logic
type OpCode int64

const (
	OpCodeUpdateEntities = 1
)

// Match is the object registered
// as a runtime.Match interface
type Match struct{}

// MatchState holds information that is passed between
// Nakama match methods
type MatchState struct {
	presences     map[string]runtime.Presence
	entityManager *ecs.EntityManager
}

// GetPrecenseList returns an array of current precenes in an array
func (state *MatchState) GetPrecenseList() []runtime.Presence {
	precenseList := []runtime.Presence{}
	for _, precense := range state.presences {
		precenseList = append(precenseList, precense)
	}
	return precenseList
}

// MatchInit is called when a new match is created
func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	state := &MatchState{
		presences:     map[string]runtime.Presence{},
		entityManager: ecs.CreateECS(),
	}
	tickRate := 5
	label := "{\"name\": \"Game World\"}"

	logger.Debug("Match Init")

	state.entityManager.AddEntity()
	state.entityManager.AddEntity()
	state.entityManager.AddEntity()

	return state, tickRate, label
}

// MatchJoinAttempt is called when a player tried to join a match
// and validates their attempt
func (m *Match) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	mState, ok := state.(*MatchState)
	if !ok {
		logger.Error("Invalid match state on join attempt!")
		return state, false, "Invalid match state!"
	}

	// Validate user is not already connected
	if _, ok := mState.presences[presence.GetUserId()]; ok {
		return mState, false, "User already logged in."
	} else {
		return mState, true, ""
	}

}

// MatchJoin is called when a player successfully joins the match
func (m *Match) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, ok := state.(*MatchState)
	if !ok {
		logger.Error("Invalid match state on join!")
		return state
	}

	for _, precense := range presences {

		// Add presence to map
		mState.presences[precense.GetUserId()] = precense
	}

	// for _, precense := range presences {
	// }

	return mState
}

// MatchLeave is called when a player leaves the match
func (m *Match) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, ok := state.(*MatchState)
	if !ok {
		logger.Error("Invalid match state on leave!")
		return state
	}
	// for _, presence := range presences {
	// }
	return mState
}

// MatchLoop is code that is executed every tick
func (m *Match) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	mState, ok := state.(*MatchState)
	if !ok {
		logger.Error("Invalid match state on match loop!")
		return state
	}

	entityData, err := mState.entityManager.GetEntitiesData()

	//	logger.Info("Data %v", string(entityData))

	if err != nil {
		logger.Error("Error getting entities data %e", err)
	} else {
		if sendErr := dispatcher.BroadcastMessage(OpCodeUpdateEntities, entityData, mState.GetPrecenseList(), nil, true); sendErr != nil {
			logger.Error(sendErr.Error())
		}
	}

	mState.entityManager.Update()

	// for _, message := range messages {

	// }

	return mState
}

// MatchTerminate is code that is executed when the match ends
func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	return state
}
