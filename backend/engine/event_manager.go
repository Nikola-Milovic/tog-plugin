package engine

import "container/heap"

type EventManager struct {
	EventPriorityQueue eventPriorityQueue
}

func CreateEventManager() *EventManager {
	em := EventManager{}
	events := make(eventPriorityQueue, 0, 100)
	em.EventPriorityQueue = events
	heap.Init(&events)
	return &em
}

func (em *EventManager) SendEvent(event Event) {
	em.EventPriorityQueue.Push(event)
}
