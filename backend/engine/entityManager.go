package engine

import (
	"fmt"
	"reflect"
)

//EntityManager is the base of the e, it holds all the Components (Structs) tightly packed in memory, it holds Actions to be handled. It also keeps a reference
// to the last active entity index
type EntityManager struct {
	maxEntities      int
	lastActiveEntity int
	ObjectPool       *ObjectPool
	Entities         []Entity
	Handlers         map[string]EventHandler
	EventManager     *EventManager
	Systems          []System
	TempSystems      map[string]TempSystem
	TempSystemMaker  map[string]func() TempSystem
	IndexMap         map[string]int //holds the indexes of entities, with their id's as keys
	AIRegistry       map[string]func() AI
	ComponentMaker   ComponentMaker
}

//CreateEntityManager creates an EntityManager, needs some more configuration, just for testing atm
func CreateEntityManager(maxSize int) *EntityManager {
	e := &EntityManager{
		Entities:         make([]Entity, 0, maxSize),
		maxEntities:      maxSize,
		lastActiveEntity: 0,
		Handlers:         make(map[string]EventHandler, 10),
		TempSystems:      make(map[string]TempSystem, 10),
		TempSystemMaker:  make(map[string]func() TempSystem, 10),
		Systems:          make([]System, 0, 10),
		IndexMap:         make(map[string]int, maxSize),
		AIRegistry:       make(map[string]func() AI, 10),
	}
	//	e.resizeComponents()

	return e
}

//Update is called every Tick of the GameLoop.
func (e *EntityManager) Update() {
	for _, sys := range e.Systems {
		sys.Update()
	}

	for _, sys := range e.TempSystems {
		sys.Update()
	}

	for _, ent := range e.Entities {
		if ent.Active {
			e.ObjectPool.AI[ent.UnitID].PerformAI(ent.Index)
		}
	}

	for len(e.EventManager.eventPriorityQueue) != 0 {
		event := e.EventManager.eventPriorityQueue.Pop().(Event)
		e.Handlers[event.ID].HandleEvent(event)
	}
}

func (e *EntityManager) AddEntity(ent NewEntityData, tag int, startOfMatch bool) (int, string) {
	data, ok := ent.Data.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("Add Entity didn't receive a NewEntityData but rather %v", reflect.TypeOf(data)))
	}

	unitID := ent.ID
	index := e.lastActiveEntity

	id, _ := shortid.Generate()

	e.IndexMap[id] = index

	e.Entities = append(e.Entities, Entity{Index: index, UnitID: unitID, Active: true, PlayerTag: ent.PlayerTag, ID: id})

	ai, ok := e.AIRegistry[unitID]

	if ok {
		_, ok := e.ObjectPool.AI[unitID]
		if !ok {
			e.ObjectPool.AI[unitID] = ai()
		}
	}

	additionalData := make(map[string]interface{}, 3)
	additionalData["position"] = ent.Position
	additionalData["tag"] = tag
	additionalData["start"] = startOfMatch
	e.ComponentMaker.AddComponents(ent.Data.(map[string]interface{}), id, additionalData)

	e.lastActiveEntity++
	//	e.resizeComponents()

	return index, id
}

//RemoveEntity swaps last entity with the deleted one, deletes entry from the IndexMap, and changes the index for the last Entity
// also calls remoteAt from objectPool to remove all components linked to the entity
func (e *EntityManager) RemoveEntity(index int) {
	entID := e.Entities[index].ID
	delete(e.IndexMap, e.Entities[index].ID)
	e.Entities[index] = e.Entities[len(e.Entities)-1]
	e.Entities = e.Entities[:len(e.Entities)-1]

	//Update the indexmap
	if index != len(e.Entities) {
		e.IndexMap[e.Entities[index].ID] = index
		e.Entities[index].Index = index
	}
	e.lastActiveEntity--
	e.ObjectPool.removeAt(index, entID)
}
func (e *EntityManager) RegisterHandler(event string, handler EventHandler) {
	if _, ok := e.Handlers[event]; ok {
		panic(fmt.Sprintf("Handler for this type of action %v is already registered", event))
	}
	fmt.Printf("Registered %v\n", event)
	e.Handlers[event] = handler
}

func (e *EntityManager) RegisterSystem(system System) {
	e.Systems = append(e.Systems, system)
}

//RegisterTempSystem
func (e *EntityManager) AddTempSystem(sysName string, data map[string]interface{}) {
	_, ok := e.TempSystems[sysName]
	if !ok {
		e.TempSystems[sysName] = e.TempSystemMaker[sysName]()
	}

	e.TempSystems[sysName].AddData(data)
}

func (e *EntityManager) RegisterTempSystem(sysName string, maker func() TempSystem) {
	e.TempSystemMaker[sysName] = maker
}

func (e *EntityManager) RemoveTempSystem(sysName string) {
	_, ok := e.TempSystems[sysName]
	if ok {
		e.TempSystems[sysName] = nil
		delete(e.TempSystems, sysName)
	}
}

func (e *EntityManager) RegisterAIMaker(unitID string, aiMaker func() AI) {
	if _, ok := e.AIRegistry[unitID]; ok {
		panic(fmt.Sprintf("AiMaker for this unit AI %v is already registered", unitID))
	}
	e.AIRegistry[unitID] = aiMaker
}
