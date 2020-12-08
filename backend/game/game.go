package game

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/heroiclabs/nakama-common/runtime"
)

// OpCode represents a enum for valid OpCodes
// used by the match logic
type OpCode int64

const TICK_RATE = 5

const (
	//OpCodeUpdateEntities is used to indicate that this data is sending the current state of the game to the clients
	OpCodeUpdateEntities = 1
)

// Match is the object registered
// as a runtime.Match interface
type Match struct{}

// MatchState holds information that is passed between
// Nakama match methods
type MatchState struct {
	presences map[string]runtime.Presence
	World     *World
}

// GetPrecenseList returns an array of current precenes in an array
func (state *MatchState) GetPrecenseList() []runtime.Presence {
	precenseList := []runtime.Presence{}
	for _, precense := range state.presences {
		precenseList = append(precenseList, precense)
	}
	return precenseList
}
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}
func testGame(w *World, logger runtime.Logger) {

	path := "/nakama/data/units.json"

	jsonData, _ := ioutil.ReadFile(path)
	var data []interface{}

	ex, err := Exists(path)
	if err != nil {
		logger.Error("File doesn't exist: %e", err.Error())
		return
	}

	logger.Info("File exists %v", ex)

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		logger.Error("Couldn't unmarshal json: %e", err.Error())
		return
	}
	logger.Debug("Unit data is %v", data)

	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])
	w.EntityManager.AddEntity(data[0])

	w.EntityManager.Entities[0].PlayerTag = 1
	w.EntityManager.Entities[1].PlayerTag = 1
	w.EntityManager.Entities[2].PlayerTag = 0
	w.EntityManager.Entities[3].PlayerTag = 0

	p1 := w.ObjectPool.Components["PositionComponent"][0].(PositionComponent)
	p2 := w.ObjectPool.Components["PositionComponent"][1].(PositionComponent)
	p3 := w.ObjectPool.Components["PositionComponent"][2].(PositionComponent)
	p4 := w.ObjectPool.Components["PositionComponent"][3].(PositionComponent)

	p1.Position = engine.Vector{0, 0}
	p2.Position = engine.Vector{2, 0}
	p3.Position = engine.Vector{0, 7}
	p4.Position = engine.Vector{4, 6}

	w.ObjectPool.Components["PositionComponent"][0] = p1
	w.ObjectPool.Components["PositionComponent"][1] = p2
	w.ObjectPool.Components["PositionComponent"][2] = p3
	w.ObjectPool.Components["PositionComponent"][3] = p4
}

// MatchInit is called when a new match is created
func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	state := &MatchState{
		presences: map[string]runtime.Presence{},
		World:     CreateWorld(),
	}
	tickRate := TICK_RATE
	label := "{\"name\": \"Game World\"}"

	testGame(state.World, logger)

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
	}
	return mState, true, ""

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

	mState.World.Counter++

	mState.World.EntityManager.Update()

	entityData, err := GetEntitiesData(mState.World)

	if err != nil {
		logger.Error("Error getting entities data %e", err.Error())
	} else {
		if sendErr := dispatcher.BroadcastMessage(OpCodeUpdateEntities, entityData, mState.GetPrecenseList(), nil, true); sendErr != nil {
			logger.Error(sendErr.Error())
		}
	}

	// for _, message := range messages {

	// }

	return mState
}

// MatchTerminate is code that is executed when the match ends
func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	return state
}
