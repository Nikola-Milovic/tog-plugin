package engine

import "container/heap"

type EventManager struct {
	eventPriorityQueue eventPriorityQueue
}

func CreateEventManager() *EventManager {
	em := EventManager{}
	events := make(eventPriorityQueue, 0, 100)
	em.eventPriorityQueue = events
	heap.Init(&events)
	return &em
}

func (em *EventManager) SendEvent(event Event) {
	em.eventPriorityQueue.Push(event)
}
