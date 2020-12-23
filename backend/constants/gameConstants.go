package constants

type UnitState string

const (
	AttackState UnitState = "Attack"
	WalkState   UnitState = "Walk"
	StunState   UnitState = "Stun"
)

const (
	MovementSpeedSlow int = 8
	MovementSpeedFast int = 4
)

//Events
const (
	MovementEvent    string = "MovementEvent"
	AttackEvent             = "AttackEvent"
	TakeDamageEvent         = "TakeDamageEvent"
	ApplyEffectEvent        = "ApplyEffectEvent"
)

//Event priorities
const (
	MovementEventPriority    int = 1
	AttackEventPriority          = 2
	TakeDamageEventPriority      = 3
	ApplyEffectEventPriority     = 4
)
