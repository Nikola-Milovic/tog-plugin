package ecs

type Handler interface {
	Handle(indx uint8)
}

type Handlers []Handler
