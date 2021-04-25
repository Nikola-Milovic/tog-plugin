package constants

const (
	MovementSpeedSlow   float32 = 4
	MovementSpeedMedium float32 = 8
	MovementSpeedFast   float32 = 12
)

const (
	TickRate = 5
)

const (
	MapWidth  = 800
	MapHeight = 512
	TileSize  = 4
	QuadrantSize = 32
)

//Sizes
const (
	StandardSize = "16"
	SmallSize    = "10"
)

const (
	StateWalking   = "w"
	StateEngaging  = "e"
	StateAttacking = "a"
	StateThinking  = "t"
)

const (
	BoidTileSize             = 15
	ImapTypeProximity string = "Proximity"
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

var IsDebug = false
