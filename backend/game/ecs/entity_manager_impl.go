package ecs

import (
	"fmt"
	"reflect"

	"github.com/Nikola-Milovic/tog-plugin/engine"
)

//EntityManager is the base of the e, it holds all the Components (Structs) tightly packed in memory, it holds Actions to be handled. It also keeps a reference
// to the last active entity index
type EntityManager struct {
	maxEntities      int
	lastActiveEntity int
	ObjectPool       *engine.ObjectPool
	Entities         []engine.Entity
	Handlers         map[string]engine.EventHandler
	EventManager     *engine.EventManager
	Systems          []engine.System
	TempSystems      map[string]engine.TempSystem
	TempSystemMaker  map[string]func(interface{}) engine.TempSystem
	IndexMap         map[string]int //holds the indexes of entities, with their id's as keys
	AIRegistry       map[string]func() engine.AI
	ComponentMaker   engine.ComponentMaker
}

//CreateEntityManager creates an EntityManager, needs some more configuration, just for testing atm
func CreateEntityManager(maxSize int) engine.EntityManagerI {
	e := &EntityManager{
		Entities:         make([]engine.Entity, 0, maxSize),
		maxEntities:      maxSize,
		lastActiveEntity: 0,
		Handlers:         make(map[string]engine.EventHandler, 10),
		TempSystems:      make(map[string]engine.TempSystem, 10),
		TempSystemMaker:  make(map[string]func(interface{}) engine.TempSystem, 10),
		Systems:          make([]engine.System, 0, 10),
		IndexMap:         make(map[string]int, maxSize),
		AIRegistry:       make(map[string]func() engine.AI, 10),
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

	for len(e.EventManager.EventPriorityQueue) != 0 {
		event := e.EventManager.EventPriorityQueue.Pop().(engine.Event)
		e.Handlers[event.ID].HandleEvent(event)
	}
}

func (e *EntityManager) AddEntity(ent engine.NewEntityData, tag int, startOfMatch bool) (int, string) {
	data, ok := ent.Data.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("Add Entity didn't receive a NewEntityData but rather %v", reflect.TypeOf(data)))
	}

	unitName := ent.ID
	index := e.lastActiveEntity

	id, _ := engine.ShortID.Generate()

	e.IndexMap[id] = index

	e.Entities = append(e.Entities, engine.Entity{Index: index, UnitID: unitName, Active: true, PlayerTag: ent.PlayerTag, ID: id})

	ai, ok := e.AIRegistry[unitName]

	if ok {
		_, ok := e.ObjectPool.AI[unitName]
		if !ok {
			e.ObjectPool.AI[unitName] = ai()
		}
	}

	additionalData := make(map[string]interface{}, 3)
	additionalData["position"] = ent.Position
	additionalData["tag"] = tag
	additionalData["start"] = startOfMatch
	additionalData["entity_id"] = id
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
	e.ObjectPool.RemoveAt(index, entID)
}
func (e *EntityManager) RegisterHandler(event string, handler engine.EventHandler) {
	if _, ok := e.Handlers[event]; ok {
		panic(fmt.Sprintf("Handler for this type of action %v is already registered", event))
	}
	fmt.Printf("Registered %v\n", event)
	e.Handlers[event] = handler
}

func (e *EntityManager) RegisterSystem(system engine.System) {
	e.Systems = append(e.Systems, system)
}

//RegisterTempSystem
func (e *EntityManager) AddTempSystem(sysName string, data map[string]interface{}, world engine.WorldI) {
	_, ok := e.TempSystems[sysName]
	if !ok {
		e.TempSystems[sysName] = e.TempSystemMaker[sysName](world)
	}

	e.TempSystems[sysName].AddData(data)
}

func (e *EntityManager) RegisterTempSystem(sysName string, maker func(interface{}) engine.TempSystem) {
	e.TempSystemMaker[sysName] = maker
}

func (e *EntityManager) RemoveTempSystem(sysName string) {
	_, ok := e.TempSystems[sysName]
	if ok {
		e.TempSystems[sysName] = nil
		delete(e.TempSystems, sysName)
	}
}

func (e *EntityManager) RegisterAIMaker(unitID string, aiMaker func() engine.AI) {
	if _, ok := e.AIRegistry[unitID]; ok {
		panic(fmt.Sprintf("AiMaker for this unit AI %v is already registered", unitID))
	}
	e.AIRegistry[unitID] = aiMaker
}

func (e *EntityManager) RegisterComponentMaker(componentName string, cm engine.ComponentMakerFun) {
	e.ComponentMaker.RegisterComponentMaker(componentName, cm)
}

func (e *EntityManager) SetComponentMaker(cm engine.ComponentMaker) {
	e.ComponentMaker = cm
}

func (e *EntityManager) SetEventManager(em *engine.EventManager) {
	e.EventManager = em
}

func (e *EntityManager) SetObjectPool(op *engine.ObjectPool) {
	e.ObjectPool = op
}

func (e *EntityManager) GetEntities() []engine.Entity {
	return e.Entities
}

func (e *EntityManager) GetIndexMap() map[string]int {
	return e.IndexMap
}

func (e *EntityManager) GetSystems() []engine.System {
	return e.Systems
}
