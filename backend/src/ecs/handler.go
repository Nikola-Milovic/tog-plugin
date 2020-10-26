package ecs

type Handler interface {
	HandleAction(indx index)
}
