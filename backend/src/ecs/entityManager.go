package ecs

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Nikola-Milovic/tog-plugin/src/action"
	//	"github.com/Nikola-Milovic/tog-plugin/src/ai"
	"github.com/Nikola-Milovic/tog-plugin/src/constants"
)

//EntityManager is the base of the e, it holds all the Components (Structs) tightly packed in memory, it holds Actions to be handled. It also keeps a reference
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
	Entities           []Entity
	movementHandler    MovementHandler
	attackHandler      AttackHandler
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
		attackHandler:      AttackHandler{},
	}

	e.movementHandler.manager = e
	e.attackHandler.manager = e
	e.resizeComponents()

	return e
}

//Update is called every Tick of the GameLoop. Here all the logic happens
// 1) we go through all of the AI's and add each action that the AI calculates into the action slice
// 2) we sort the actions so we use the same Handlers in consecutive fashion to maximize CPU Cache, ie. 10 Movement Actions will all use the same Position slice which will already be loaded in Cache
// 3) dispatch Actions to corresponding Handlers
func (e *EntityManager) Update() {
	for index, ai := range e.AIComponents {
		e.Actions[index] = ai.AI.CalculateAction(index, e)
	}

	e.sortActions()

	//	fmt.Printf("Entities %v, Actions %v, last active %v \n", len(e.AIComponents), len(e.Actions), e.lastActiveEntity)

	for index, act := range e.Actions {
		switch act.(type) {
		case action.EmptyAction: // EmptyActions are used for entities who aren't doing anything and they will always be last in the slice, so when we encounter the first one, break
			break
		case action.MovementAction:
			e.movementHandler.HandleAction(index)
		case action.AttackAction:
			e.attackHandler.HandleAction(index)
		}
	}
}

//AddEntity adds an entity and all of its components to the Manager, WIP
func (e *EntityManager) AddEntity() {

	fmt.Println("Add Entity")

	ai := KnightAI{}

	e.Entities = append(e.Entities, Entity{PlayerTag: 1, Index: e.lastActiveEntity, Size: constants.V2{X: 16, Y: 16}})
	e.AIComponents = append(e.AIComponents, AIComponent{AI: ai})
	e.PositionComponents = append(e.PositionComponents, PositionComponent{Position: constants.V2{X: 20, Y: 20}})
	e.AttackComponents = append(e.AttackComponents, AttackComponent{Type: "phys", Range: 10})
	e.MovementComponents = append(e.MovementComponents, MovementComponent{Speed: 5})
	e.lastActiveEntity++

	ai3 := KnightAI{}
	e.Entities = append(e.Entities, Entity{PlayerTag: 1, Index: e.lastActiveEntity, Size: constants.V2{X: 16, Y: 16}})
	e.AIComponents = append(e.AIComponents, AIComponent{AI: ai3})
	e.PositionComponents = append(e.PositionComponents, PositionComponent{Position: constants.V2{X: 20, Y: 40}})
	e.AttackComponents = append(e.AttackComponents, AttackComponent{Type: "phys", Range: 10})
	e.MovementComponents = append(e.MovementComponents, MovementComponent{Speed: 5})
	e.lastActiveEntity++

	ai2 := KnightAI{}
	e.Entities = append(e.Entities, Entity{PlayerTag: 0, Index: e.lastActiveEntity, Size: constants.V2{X: 16, Y: 16}})
	e.AIComponents = append(e.AIComponents, AIComponent{AI: ai2})
	e.PositionComponents = append(e.PositionComponents, PositionComponent{Position: constants.V2{X: 200, Y: 300}})
	e.AttackComponents = append(e.AttackComponents, AttackComponent{Type: "phys", Range: 10})
	e.MovementComponents = append(e.MovementComponents, MovementComponent{Speed: 5})
	e.lastActiveEntity++

	e.Entities = append(e.Entities, Entity{PlayerTag: 1, Index: e.lastActiveEntity, Size: constants.V2{X: 16, Y: 16}})
	e.AIComponents = append(e.AIComponents, AIComponent{AI: ai2})
	e.PositionComponents = append(e.PositionComponents, PositionComponent{Position: constants.V2{X: 20, Y: 100}})
	e.AttackComponents = append(e.AttackComponents, AttackComponent{Type: "phys", Range: 10})
	e.MovementComponents = append(e.MovementComponents, MovementComponent{Speed: 5})
	e.lastActiveEntity++

	e.Entities = append(e.Entities, Entity{PlayerTag: 0, Index: e.lastActiveEntity, Size: constants.V2{X: 16, Y: 16}})
	e.AIComponents = append(e.AIComponents, AIComponent{AI: ai2})
	e.PositionComponents = append(e.PositionComponents, PositionComponent{Position: constants.V2{X: 300, Y: 300}})
	e.AttackComponents = append(e.AttackComponents, AttackComponent{Type: "phys", Range: 10})
	e.MovementComponents = append(e.MovementComponents, MovementComponent{Speed: 5})
	e.lastActiveEntity++

	e.Actions = e.Actions[:e.lastActiveEntity]

}

//RemoveEntity WIP
func (e *EntityManager) RemoveEntity() {

}

// Used to resize all of the component slices down to size of active entities, so we don't waste loops in the Update
func (e *EntityManager) resizeComponents() {
	e.AttackComponents = e.AttackComponents[:e.lastActiveEntity]
	e.MovementComponents = e.MovementComponents[:e.lastActiveEntity]
	e.PositionComponents = e.PositionComponents[:e.lastActiveEntity]
	e.AIComponents = e.AIComponents[:e.lastActiveEntity]
	e.Actions = e.Actions[:e.lastActiveEntity]
	e.Entities = e.Entities[:e.lastActiveEntity]
}

//sortActions is used to sort the Actions based on priority, this provides us with grouping of Actions which will simplify CPU Caching
func (e *EntityManager) sortActions() {
	sort.Sort(sortByActionPriority(e.Actions))
}

//GetEntitiesData gets the data of all entities and packs them into []byte, used to send the clients necessary data to reconstruct the current state of the game
//TODO: add batching instead of sending all the data at once
func (e *EntityManager) GetEntitiesData() ([]byte, error) {
	size := e.lastActiveEntity
	entities := make([]EntityData, 0, size+1)

	for i := 0; i < size; i++ {
		//	fmt.Printf("I at %v am at position %v \n", i, e.PositionComponents[i].Position)
		entities = append(entities, EntityData{
			Index:    i,
			Position: e.PositionComponents[i].Position,
			State:    e.Actions[i].GetActionState(),
			Velocity: e.MovementComponents[i].Velocity,
		})
	}

	data, err := json.Marshal(&entities)
	return data, err
}

func (e *EntityManager) getNearbyEntities(maxDistance int, position constants.V2, index int) []int {
	nearbyEntities := make([]int, 0, len(e.Entities))

	for idx, posComp := range e.PositionComponents {
		if idx == index {
			continue
		}
		dist := position.Distance(posComp.Position)
		if dist <= maxDistance {
			//	fmt.Printf("Found entity at %v, distance to %v \n", idx, dist)
			nearbyEntities = append(nearbyEntities, idx)
		}
	}

	return nearbyEntities

}

//Used to sort actions by priority so we will save memory with CPU caching as the actions will be of the same type
type sortByActionPriority []action.Action

func (a sortByActionPriority) Len() int           { return len(a) }
func (a sortByActionPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByActionPriority) Less(i, j int) bool { return a[i].GetPriority() < a[j].GetPriority() }
