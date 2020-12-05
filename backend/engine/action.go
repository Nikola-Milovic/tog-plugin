package engine

//Action represents a single activity/ action that the entity will perform that Tick, AI decides on the best action that the
//given entity should perform
type Action interface {
	GetPriority() int
	GetActionType() string
}
