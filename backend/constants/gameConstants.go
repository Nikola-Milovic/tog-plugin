package constants

const (
	MovementSpeedSlow int = 8
	MovementSpeedFast int = 4
)

//Events
const (
	MovementEvent            string = "MovementEvent"
	AttackEvent                     = "AttackEvent"
	TakeDamageEvent                 = "TakeDamageEvent"
	ApplyEffectEvent                = "ApplyEffectEvent"
	AbilityCastEvent                = "AbilityCastEvent"
	SingleTargetAbilityEvent        = "SingleTargetAbilityEvent"
	SummonAbilityEvent              = "SummonAbilityEvent"
	LineShotAbilityEvent            = "LineShotAbilityEvent"
)

//Event priorities
const (
	MovementEventPriority            int = 1
	AttackEventPriority                  = 2
	TakeDamageEventPriority              = 3
	ApplyEffectEventPriority             = 4
	AbilityCastEventPriority             = 5
	SingleTargetAbilityEventPriority     = 6
	SummonAbilityEventPriority           = 7
	LineShotEventPriority                = 8
)
