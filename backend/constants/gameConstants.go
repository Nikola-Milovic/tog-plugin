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
	MovementEvent   string = "MovementEvent"
	AttackEvent            = "AttackEvent"
	TakeDamageEvent        = "TakeDamageEvent"
)
