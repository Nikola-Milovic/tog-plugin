package ecs

type Handler interface {
	HandleAction(indx int)
}
