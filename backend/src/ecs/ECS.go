package ecs

import (
	"backend/src/action"
	"sort"
)

type ECS struct {
	maxEntities      uint8
	indexPool        []uint8
	lastActiveEntity uint8

	handlers Handlers

	Actions []action.Action

	AttackComponents   []AttackComponent
	MovementComponents []MovementComponent
	PositionComponents []PositionComponent

	Entities []Entity
}

func (ecs *ECS) Update() {
	ecs.sortActions()

	for index, act := range ecs.Actions {
		switch a := act.(type) {
		case action.MovementAction:

		}
	}

}

func (ecs *ECS) AddEntity() {

}

func (ecs *ECS) RemoveEntity() {

}

//Used to sort actions by priority so we will save memory with CPU caching as the actions will be of the same type
type sortByActionPriority []action.Action

func (a sortByActionPriority) Len() int           { return len(a) }
func (a sortByActionPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByActionPriority) Less(i, j int) bool { return a[i].GetPriority() < a[j].GetPriority() }

func (esc *ECS) sortActions() {
	sort.Sort(sortByActionPriority(esc.Actions))
}
