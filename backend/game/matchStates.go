package game

type MatchState int

const (
	MatchStartedState          MatchState = 0
	MatchPreperationState      MatchState = 1
	MatchWaitingForPlayerState MatchState = 2
)

type MatchStateMessage struct {
	MatchState MatchState `json:"match_state"`
}
