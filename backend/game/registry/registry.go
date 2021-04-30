package registry

import (
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	ai "github.com/Nikola-Milovic/tog-plugin/game/AI"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/ecs"
	"github.com/Nikola-Milovic/tog-plugin/game/grid"
	"github.com/Nikola-Milovic/tog-plugin/game/handlers"
	"github.com/Nikola-Milovic/tog-plugin/game/systems"
	"github.com/Nikola-Milovic/tog-plugin/game/tempsys"
)

func RegisterWorld(w *game.World) {
	w.SetupECS(ecs.CreateEntityManager(30))
	w.Grid = grid.CreateGrid(w)

	registerAIMakers(w)
	registerComponentMakers(w)
	registerHandlers(w)
	registerSystems(w)
	registerTempSystems(w)
}

func registerComponentMakers(w *game.World) {
	w.EntityManager.SetComponentMaker(components.CreateComponentMaker(w))
	w.EntityManager.RegisterComponentMaker("MovementComponent", components.MovementComponentMaker)
	w.EntityManager.RegisterComponentMaker("PositionComponent", components.PositionComponentMaker)
	w.EntityManager.RegisterComponentMaker("AttackComponent", components.AttackComponentMaker)
	w.EntityManager.RegisterComponentMaker("StatsComponent", components.StatsComponentMaker)
	w.EntityManager.RegisterComponentMaker("EffectsComponent", components.EffectsComponentMaker)
	w.EntityManager.RegisterComponentMaker("FlagComponent", components.FlagComponentMaker)
	w.EntityManager.RegisterComponentMaker("AbilitiesComponent", components.AbilitiesComponentMaker)
}

func registerHandlers(w *game.World) {
	w.EntityManager.RegisterHandler(constants.TakeDamageEvent, handlers.TakeDamageEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.ApplyEffectEvent, handlers.ApplyEffectEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.TriggerActionEvent, handlers.ActionHandler{World: w})
	//Abilities
	w.EntityManager.RegisterHandler(constants.AbilityCastEvent, handlers.AbilityCastEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.SingleTargetAbilityEvent, handlers.SingleTargetAbilityEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.SummonAbilityEvent, handlers.SummonAbilityEventHandler{World: w})
	w.EntityManager.RegisterHandler(constants.LineShotAbilityEvent, handlers.LineshotAbilityEventHandler{World: w})
}

func registerAIMakers(w *game.World) { // TODO move AI makers as constants
	w.EntityManager.RegisterAIMaker("knight", func() engine.AI { return ai.GenericAI{World: w} })
	w.EntityManager.RegisterAIMaker("archer", func() engine.AI { return ai.GenericAI{World: w} })
	w.EntityManager.RegisterAIMaker("gob_beast_master", func() engine.AI { return ai.GoblinBeastMasterAI{World: w} })
	w.EntityManager.RegisterAIMaker("gob_spear", func() engine.AI { return ai.GenericAI{World: w} })
	w.EntityManager.RegisterAIMaker("s_wolf", func() engine.AI { return ai.GenericAI{World: w} })
}

func registerSystems(w *game.World) {
	w.EntityManager.RegisterSystem(systems.DeathSystem{World: w})

	w.EntityManager.RegisterSystem(systems.DotSystem{World: w})
	w.EntityManager.RegisterSystem(systems.DurationSystem{World: w})

	//Movement and logic
	w.EntityManager.RegisterSystem(systems.MovementSystem{World: w})
	w.EntityManager.RegisterSystem(systems.TargetingSystem{World: w})
	w.EntityManager.RegisterSystem(systems.CollisionSystem{World: w})
	w.EntityManager.RegisterSystem(systems.AttackSystem{World: w})

	w.EntityManager.RegisterSystem(systems.ClientSystem{World: w})
}

func registerTempSystems(w *game.World) {
	w.EntityManager.RegisterTempSystem("LineshotTempSystem", tempsys.CreateLineShotTempSystem)
}
