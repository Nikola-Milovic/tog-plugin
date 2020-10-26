package ecs

type Handler interface {
	Handle(indx uint16)
}
