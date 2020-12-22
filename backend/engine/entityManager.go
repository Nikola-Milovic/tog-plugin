package engine

import (
	"fmt"
	"reflect"
)

//EntityManager is the base of the e, it holds all the Components (Structs) tightly packed in memory, it holds Actions to be handled. It also keeps a reference
// to the last active entity index
type EntityManager struct {
	maxEntities       int
	lastActiveEntity  int
	ObjectPool        *ObjectPool
	ComponentRegistry map[string]ComponentMaker
	AIRegistry        map[string]func() AI
	Entities          []Entity
	Handlers          map[string]EventHandler
	EventManager      *EventManager
}

//CreateEntityManager creates an EntityManager, needs some more configuration, just for testing atm
func CreateEntityManager(maxSize int) *EntityManager {
	e := &EntityManager{
		Entities:          make([]Entity, 0, maxSize),
		maxEntities:       maxSize,
		lastActiveEntity:  0,
		ComponentRegistry: make(map[string]ComponentMaker, 10),
		Handlers:          make(map[string]EventHandler, 10),
		AIRegistry:        make(map[string]func() AI, 10),
	}

	//	e.resizeComponents()

	return e
}

func (e *EntityManager) RegisterComponentMaker(componentName string, maker ComponentMaker) {
	if _, ok := e.ComponentRegistry[componentName]; ok {
		panic(fmt.Sprintf("Component maker for component %v is already registered", componentName))
	}
	e.ComponentRegistry[componentName] = maker
}

func (e *EntityManager) RegisterHandler(event string, handler EventHandler) {
	if _, ok := e.Handlers[event]; ok {
		panic(fmt.Sprintf("Handler for this type of action %v is already registered", event))
	}
	fmt.Printf("Registered %v\n", event)
	e.Handlers[event] = handler
}

func (e *EntityManager) RegisterAIMaker(unitID string, aiMaker func() AI) {
	if _, ok := e.AIRegistry[unitID]; ok {
		panic(fmt.Sprintf("AiMaker for this unit AI %v is already registered", unitID))
	}
	e.AIRegistry[unitID] = aiMaker
}

//Update is called every Tick of the GameLoop.
func (e *EntityManager) Update() {
	for _, ent := range e.Entities {
		if ent.Active {
			e.ObjectPool.AI[ent.ID].PerformAI(ent.Index)
		}
	}

	for len(e.EventManager.eventPriorityQueue) != 0 {
		event := e.EventManager.eventPriorityQueue.Pop().(Event)
		e.Handlers[event.ID].HandleEvent(event)
	}
}

//AddEntity adds an entity and all of its components to the object pool
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
	//	e.resizeComponents()

	return index
}

//RemoveEntity WIP
func (e *EntityManager) RemoveEntity(index int) {
	e.Entities[index].Active = false
	//e.AvailableIndexPool = append(e.AvailableIndexPool, index)
	//	e.resizeComponents()
}

// // Used to resize all of the component slices down to size of active entities, so we don't waste loops in the Update
// func (e *EntityManager) resizeComponents() {
// 	e.Actions = e.Actions[:e.lastActiveEntity]
// 	e.Entities = e.Entities[:e.lastActiveEntity]
// }
