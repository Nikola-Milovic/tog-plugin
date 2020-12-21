package engine

//EventHandler interface indicates that it handles a certain action. Each handler should implement it and provide a way to resolve the action it is assigned to
//Handlers are like Systems in regular ECS but a EventHandler is a more fitting name, as it handles a specific event in game, Ie MovementEventHandler, AttackEventHandler, DeathEventHandler...
type EventHandler interface {
	HandleEvent(ev Event)
}
