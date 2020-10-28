package ecs

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Nikola-Milovic/tog-plugin/src/action"
	"github.com/Nikola-Milovic/tog-plugin/src/ai"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

//EntityManager is the base of the ECS, it holds all the Components (Structs) tightly packed in memory, it holds Actions to be handled. It also keeps a reference
// to the last active entity index
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

//CreateEntityManager creates an EntityManager, needs some more configuration, just for testing atm
func CreateEntityManager() *EntityManager {
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

//Update is called every Tick of the GameLoop. Here all the logic happens
// 1) we go through all of the AI's and add each action that the AI calculates into the action slice
// 2) we sort the actions so we use the same Handlers in consecutive fashion to maximize CPU Cache, ie. 10 Movement Actions will all use the same Position slice which will already be loaded in Cache
// 3) dispatch Actions to corresponding Handlers
func (ecs *EntityManager) Update() {
	for index, ai := range ecs.AIComponents {
		ecs.Actions[index] = ai.AI.CalculateAction(index)
	}

	ecs.sortActions()

	//	fmt.Printf("Entities %v, Actions %v, last active %v \n", len(ecs.AIComponents), len(ecs.Actions), ecs.lastActiveEntity)

	for index, act := range ecs.Actions {
		switch act.(type) {
		case action.EmptyAction: // EmptyActions are used for entities who aren't doing anything and they will always be last in the slice, so when we encounter the first one, break
			break
		case action.MovementAction:
			ecs.movementHandler.HandleAction(index)
		}
	}
}

//AddEntity adds an entity and all of its components to the Manager, WIP
func (ecs *EntityManager) AddEntity() {

	fmt.Println("Add Entity")

	ai := ai.KnightAI{}

	ecs.AIComponents = append(ecs.AIComponents, AIComponent{AI: ai})
	ecs.PositionComponents = append(ecs.PositionComponents, PositionComponent{Position: constants.V2{X: 10, Y: 10}})
	ecs.MovementComponents = append(ecs.MovementComponents, MovementComponent{Speed: 5})
	ecs.Actions = ecs.Actions[:ecs.lastActiveEntity+1]
	ecs.lastActiveEntity++
}

//RemoveEntity WIP
func (ecs *EntityManager) RemoveEntity() {

}

// Used to resize all of the component slices down to size of active entities, so we don't waste loops in the Update
func (ecs *EntityManager) resizeComponents() {
	ecs.AttackComponents = ecs.AttackComponents[:ecs.lastActiveEntity]
	ecs.MovementComponents = ecs.MovementComponents[:ecs.lastActiveEntity]
	ecs.PositionComponents = ecs.PositionComponents[:ecs.lastActiveEntity]
	ecs.AIComponents = ecs.AIComponents[:ecs.lastActiveEntity]
	ecs.Actions = ecs.Actions[:ecs.lastActiveEntity]
}

//sortActions is used to sort the Actions based on priority, this provides us with grouping of Actions which will simplify CPU Caching
func (ecs *EntityManager) sortActions() {
	sort.Sort(sortByActionPriority(ecs.Actions))
}

//GetEntitiesData gets the data of all entities and packs them into []byte, used to send the clients necessary data to reconstruct the current state of the game
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

//Used to sort actions by priority so we will save memory with CPU caching as the actions will be of the same type
type sortByActionPriority []action.Action

func (a sortByActionPriority) Len() int           { return len(a) }
func (a sortByActionPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByActionPriority) Less(i, j int) bool { return a[i].GetPriority() < a[j].GetPriority() }
