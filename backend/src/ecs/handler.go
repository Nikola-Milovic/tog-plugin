package ecs

//Handler interface indicates that it handles a certain action. Each handler should implement it and provide a way to resolve the action it is assigned to
//Handlers are like Systems in regular ECS but a Handler is a more fitting name, as it handles a specific action in game, Ie MovementHandler, AttackHandler, DeathHandler...
type Handler interface {
	HandleAction(indx int)
}
