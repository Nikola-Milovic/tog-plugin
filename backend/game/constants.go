package game

type UnitState string

const (
	AttackState UnitState = "Attack"
	WalkState UnitState = "Walk"
	StunState UnitState = "Stun"
)