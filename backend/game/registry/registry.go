package registry

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	ai "github.com/Nikola-Milovic/tog-plugin/game/AI"
	"github.com/Nikola-Milovic/tog-plugin/game/client"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/handlers"
	"github.com/Nikola-Milovic/tog-plugin/game/systems"
)

func RegisterWorld(w *game.World) {
	registerAIMakers(w)
	registerComponentMakers(w)
	registerHandlers(w)
	registerSystems(w)
	registerEventConverters(w)
}

func registerComponentMakers(w *game.World) {
	w.EntityManager.RegisterComponentMaker("MovementComponent", components.MovementComponentMaker)
	w.EntityManager.RegisterComponentMaker("PositionComponent", components.PositionComponentMaker)
	w.EntityManager.RegisterComponentMaker("AttackComponent", components.AttackComponentMaker)
	w.EntityManager.RegisterComponentMaker("StatsComponent", components.StatsComponentMaker)
	w.EntityManager.RegisterComponentMaker("EffectsComponent", components.EffectsComponentMaker)
	w.EntityManager.RegisterComponentMaker("AbilitiesComponent", components.AbilitiesComponentMaker)
}

func registerHandlers(w *game.World) {
	w.EntityManager.RegisterHandler(constants.MovementEvent, handlers.MovementEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.AttackEvent, handlers.AttackEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.TakeDamageEvent, handlers.TakeDamageEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.ApplyEffectEvent, handlers.ApplyEffectEventHandler{World: w})
	//Abilities
	w.EntityManager.RegisterHandler(constants.AbilityCastEvent, handlers.AbilityCastEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.SingleTargetAbilityEvent, handlers.SingleTargetAbilityEventHandler{World: w})
}

func registerAIMakers(w *game.World) {
	w.EntityManager.RegisterAIMaker("knight", func() engine.AI { return ai.KnightAI{World: w} })
	w.EntityManager.RegisterAIMaker("archer", func() engine.AI { return ai.ArcherAI{World: w} })
}

func registerSystems(w *game.World) {
	w.EntityManager.RegisterSystem(systems.DeathSystem{World: w})
	w.EntityManager.RegisterSystem(systems.DotSystem{World: w})
	w.EntityManager.RegisterSystem(systems.DurationSystem{World: w})
	w.EntityManager.RegisterSystem(systems.MovementSystem{World: w})
	w.EntityManager.RegisterSystem(systems.AttackSystem{World: w})
}

func registerEventConverters(w *game.World) {
	w.ClientEventManager.RegisterClientEventConverter(client.AttackEventConverter, constants.AttackEvent)
}
