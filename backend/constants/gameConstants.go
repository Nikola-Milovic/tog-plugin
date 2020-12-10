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

	AttackSpeedSlow int = 8
	AttackSpeedFast int = 4
)

const (
	ActionTypeMovement string = "movement"
	ActionTypeAttack          = "attack"
	ActionTypeEmpty           = "empty"
)
