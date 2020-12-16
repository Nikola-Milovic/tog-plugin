package engine

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/Nikola-Milovic/tog-plugin/constants"
)

//EntityManager is the base of the e, it holds all the Components (Structs) tightly packed in memory, it holds Actions to be handled. It also keeps a reference
// to the last active entity index
type EntityManager struct {
	maxEntities        int
	lastActiveEntity   int
	ObjectPool         *ObjectPool
	ComponentRegistry  map[string]ComponentMaker
	AIRegistry         map[string]func() AI
	Actions            []Action
	Entities           []Entity
	Handlers           map[string]Handler
	AvailableIndexPool []int
}

//CreateEntityManager creates an EntityManager, needs some more configuration, just for testing atm
func CreateEntityManager(maxSize int) *EntityManager {
	e := &EntityManager{
		Entities:           make([]Entity, 0, maxSize),
		maxEntities:        maxSize,
		lastActiveEntity:   0,
		Actions:            make([]Action, 0, maxSize),
		ComponentRegistry:  make(map[string]ComponentMaker, 10),
		Handlers:           make(map[string]Handler, 10),
		AIRegistry:         make(map[string]func() AI, 10),
		AvailableIndexPool: make([]int, 0, maxSize),
	}

	e.resizeComponents()

	return e
}

func (e *EntityManager) RegisterComponentMaker(componentName string, maker ComponentMaker) {
	if _, ok := e.ComponentRegistry[componentName]; ok {
		panic(fmt.Sprintf("Component maker for component %v is already registered", componentName))
	}
	e.ComponentRegistry[componentName] = maker
}

func (e *EntityManager) RegisterHandler(actionType string, handler Handler) {
	if _, ok := e.Handlers[actionType]; ok {
		panic(fmt.Sprintf("Handler for this type of action %v is already registered", actionType))
	}
	e.Handlers[actionType] = handler
}

func (e *EntityManager) RegisterAIMaker(unitID string, aiMaker func() AI) {
	if _, ok := e.AIRegistry[unitID]; ok {
		panic(fmt.Sprintf("AiMaker for this unit AI %v is already registered", unitID))
	}
	e.AIRegistry[unitID] = aiMaker
}

//Update is called every Tick of the GameLoop. Here all the logic happens
// 1) we go through all of the entities and find their AI, then we add each action that the AI calculates into the action slice
// 2) we sort the actions so we use the same Handlers in consecutive fashion to maximize CPU Cache, ie. 10 Movement Actions will all use the same Position slice which will already be loaded in Cache
// 3) dispatch Actions to corresponding Handlers
func (e *EntityManager) Update() {
	for _, ent := range e.Entities {
		if ent.Active {
			e.Actions[ent.Index] = e.ObjectPool.AI[ent.ID].CalculateAction(ent.Index)
		}
	}

	e.sortActions()

	for _, act := range e.Actions {
		if act.GetActionType() == constants.ActionTypeEmpty {
			break
		}
		e.Handlers[act.GetActionType()].HandleAction(act)
	}
}

//AddEntity adds an entity and all of its components to the Manager, WIP
func (e *EntityManager) AddEntity(entityData NewEntityData) int {
	data, ok := entityData.Data.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("Add Entity didn't receive a NewEntityData but rather %v", reflect.TypeOf(data)))
	}

	components, ok := data["Components"].(map[string]interface{})
	if !ok {
		panic(fmt.Sprint("Added entity doesn't have components"))
	}

	unitID := entityData.ID
	index := e.lastActiveEntity

	//Eg key = MovementComponent, data is MovementSpeed, MovementType etc
	for key, data := range components {
		maker, ok := e.ComponentRegistry[key] // MovementComponentMaker, returns a MovementComponent
		if !ok {
			panic(fmt.Sprintf("No registered maker for the component %s", key))
		}
		component := maker(data)
		e.ObjectPool.addComponent(component)
	}

	e.Entities = append(e.Entities, Entity{Index: index, ID: unitID, Active: true, PlayerTag: entityData.PlayerTag})

	ai, ok := e.AIRegistry[unitID]

	if ok {
		_, ok := e.ObjectPool.AI[unitID]
		if !ok {
			e.ObjectPool.AI[unitID] = ai()
		}
	}

	e.lastActiveEntity++
	e.resizeComponents()

	return index
}

//RemoveEntity WIP
func (e *EntityManager) RemoveEntity(index int) {
	e.Entities[index].Active = false
	e.AvailableIndexPool = append(e.AvailableIndexPool, index)
	//	e.resizeComponents()
}

// Used to resize all of the component slices down to size of active entities, so we don't waste loops in the Update
func (e *EntityManager) resizeComponents() {
	e.Actions = e.Actions[:e.lastActiveEntity]
	e.Entities = e.Entities[:e.lastActiveEntity]
}

//sortActions is used to sort the Actions based on priority, this provides us with grouping of Actions which will simplify CPU Caching
func (e *EntityManager) sortActions() {
	sort.Sort(sortByActionPriority(e.Actions))
}

//Used to sort actions by priority so we will save memory with CPU caching as the actions will be of the same type
type sortByActionPriority []Action

func (a sortByActionPriority) Len() int           { return len(a) }
func (a sortByActionPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByActionPriority) Less(i, j int) bool { return a[i].GetPriority() > a[j].GetPriority() }
