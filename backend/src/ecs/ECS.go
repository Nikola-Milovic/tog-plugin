package ecs

import (
	"backend/src/action"
	"sort"
)

type ECS struct {
	maxEntities      uint8
	indexPool        []uint8
	lastActiveEntity uint8

	Actions []action.Action

	AttackComponents   []AttackComponent
	MovementComponents []MovementComponent
	PositionComponents []PositionComponent
	AIComponents       []AIComponent

	movementHandler MovementHandler
}

func CreateECS() *ECS {
	return &ECS{
		maxEntities:        10,
		indexPool:          make([]uint8, 10),
		lastActiveEntity:   0,
		Actions:            make([]action.Action, 10),
		AttackComponents:   make([]AttackComponent, 10),
		MovementComponents: make([]MovementComponent, 10),
		PositionComponents: make([]PositionComponent, 10),
		AIComponents:       make([]AIComponent, 10),
		movementHandler:    MovementHandler{},
	}
}

func (ecs *ECS) Update() {

	idx := uint16(0)
	for _, ai := range ecs.AIComponents {
		ecs.Actions[idx] = ai.AI.CalculateAction(idx)
		idx++
	}

	ecs.sortActions()

	idx = uint16(0)
	for _, act := range ecs.Actions {
		switch act.(type) {
		case action.MovementAction:
			ecs.movementHandler.Handle(idx)

		}

		idx++
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
