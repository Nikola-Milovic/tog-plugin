package constants

const (
	MovementSpeedSlow int = 1
	MovementSpeedMedium int = 2
	MovementSpeedFast int = 3
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
	TriggerActionEvent              = "TriggerActionEvent"
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
	TriggerActionEventPriority           = 9
)
