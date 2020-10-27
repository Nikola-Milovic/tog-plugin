package ecs

import (
	"encoding/json"
	"fmt"
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
	e := &EntityManager{
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

	e.movementHandler.manager = e

	e.resizeComponents()

	return e
}

func (ecs *EntityManager) Update() {
	for index, ai := range ecs.AIComponents {
		ecs.Actions[index] = ai.AI.CalculateAction(index)
	}

	ecs.sortActions()

	//	fmt.Printf("Entities %v, Actions %v, last active %v \n", len(ecs.AIComponents), len(ecs.Actions), ecs.lastActiveEntity)

	for index, act := range ecs.Actions {
		switch act.(type) {
		case action.MovementAction:
			ecs.movementHandler.HandleAction(index)
		}
	}
}

func (ecs *EntityManager) AddEntity() {

	fmt.Println("Add Entity")

	ai := ai.KnightAI{}

	ecs.AIComponents = append(ecs.AIComponents, AIComponent{AI: ai})
	ecs.PositionComponents = append(ecs.PositionComponents, PositionComponent{Position: constants.V2{X: 10, Y: 10}})
	ecs.MovementComponents = append(ecs.MovementComponents, MovementComponent{Speed: 5})
	ecs.Actions = ecs.Actions[:ecs.lastActiveEntity+1]
	ecs.lastActiveEntity++
}

func (ecs *EntityManager) RemoveEntity() {

}

func (ecs *EntityManager) resizeComponents() {
	ecs.AttackComponents = ecs.AttackComponents[:ecs.lastActiveEntity]
	ecs.MovementComponents = ecs.MovementComponents[:ecs.lastActiveEntity]
	ecs.PositionComponents = ecs.PositionComponents[:ecs.lastActiveEntity]
	ecs.AIComponents = ecs.AIComponents[:ecs.lastActiveEntity]
	ecs.Actions = ecs.Actions[:ecs.lastActiveEntity]
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
			Index:    i,
			Position: ecs.PositionComponents[i].Position,
			Action:   "walk",
		}
	}

	data, err := json.Marshal(&entities)
	return data, err
}
