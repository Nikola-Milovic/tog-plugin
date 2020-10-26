package ecs

import (
	"encoding/json"
	"sort"

	"github.com/Nikola-Milovic/tog-plugin/src/action"
	"github.com/Nikola-Milovic/tog-plugin/src/ai"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

type EntityManager struct {
	maxEntities      int
	indexPool        []int
	lastActiveEntity int

	Actions []action.Action

	AttackComponents   []AttackComponent
	MovementComponents []MovementComponent
	PositionComponents []PositionComponent
	AIComponents       []AIComponent

	movementHandler MovementHandler
}

func CreateECS() *EntityManager {
	return &EntityManager{
		maxEntities:        10,
		indexPool:          make([]int, 10),
		lastActiveEntity:   0,
		Actions:            make([]action.Action, 10),
		AttackComponents:   make([]AttackComponent, 10),
		MovementComponents: make([]MovementComponent, 10),
		PositionComponents: make([]PositionComponent, 10),
		AIComponents:       make([]AIComponent, 10),
		movementHandler:    MovementHandler{},
	}
}

func (ecs *EntityManager) Update() {

	for index, ai := range ecs.AIComponents {
		if ai.AI != nil {
			ecs.Actions[index] = ai.AI.CalculateAction(index)
		}
	}

	//	ecs.sortActions()

	// for index, act := range ecs.Actions {
	// 	switch act.(type) {
	// 	case action.MovementAction:
	// 		ecs.movementHandler.HandleAction(index)

	// 	}
	// }
}

func (ecs *EntityManager) AddEntity() {
	ecs.AIComponents[ecs.lastActiveEntity] = AIComponent{&ai.KnightAI{}}
	ecs.PositionComponents[ecs.lastActiveEntity] = PositionComponent{Position: constants.V2{X: 10, Y: 10}}
	ecs.MovementComponents[ecs.lastActiveEntity] = MovementComponent{Speed: 5}
	ecs.lastActiveEntity++
}

func (ecs *EntityManager) RemoveEntity() {

}

//Used to sort actions by priority so we will save memory with CPU caching as the actions will be of the same type
type sortByActionPriority []action.Action

func (a sortByActionPriority) Len() int           { return len(a) }
func (a sortByActionPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByActionPriority) Less(i, j int) bool { return a[i].GetPriority() < a[j].GetPriority() }

func (esc *EntityManager) sortActions() {
	sort.Sort(sortByActionPriority(esc.Actions))
}

func (ecs *EntityManager) GetEntitiesData() ([]byte, error) {
	size := ecs.lastActiveEntity
	entities := make([]EntityData, size)

	for i := 0; i < size; i++ {
		entities[i] = EntityData{
			Position: ecs.PositionComponents[i].Position,
			Action:   "walk",
		}
	}

	data, err := json.Marshal(&entities)
	return data, err
}
