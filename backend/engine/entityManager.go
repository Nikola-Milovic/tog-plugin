package engine

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

//EntityManager is the base of the e, it holds all the Components (Structs) tightly packed in memory, it holds Actions to be handled. It also keeps a reference
// to the last active entity index
type EntityManager struct {
	maxEntities       int
	lastActiveEntity  int
	ObjectPool        *ObjectPool
	ComponentRegistry map[string]ComponentMaker
	Actions           []Action
	Entities          []Entity
}

//CreateEntityManager creates an EntityManager, needs some more configuration, just for testing atm
func CreateEntityManager() *EntityManager {
	e := &EntityManager{
		maxEntities:       10,
		lastActiveEntity:  0,
		Actions:           make([]Action, 0, 10),
		ComponentRegistry: make(map[string]ComponentMaker, 10),
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

//Update is called every Tick of the GameLoop. Here all the logic happens
// 1) we go through all of the AI's and add each action that the AI calculates into the action slice
// 2) we sort the actions so we use the same Handlers in consecutive fashion to maximize CPU Cache, ie. 10 Movement Actions will all use the same Position slice which will already be loaded in Cache
// 3) dispatch Actions to corresponding Handlers
func (e *EntityManager) Update() {
	e.sortActions()

	for _, act := range e.Actions {
		println(act)
	}
}

//AddEntity adds an entity and all of its components to the Manager, WIP
func (e *EntityManager) AddEntity(entityData interface{}) {
	data, ok := entityData.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("Add Entity didn't receive a map[string]interface but rather %v", reflect.TypeOf(data)))
	}

	components, ok := data["Components"].(map[string]interface{})
	if !ok {
		panic(fmt.Sprint("Added entity doesn't have components"))
	}

	//Eg key = MovementComponent, data is MovementSpeed, MovementType etc
	for key, data := range components {
		maker, ok := e.ComponentRegistry[key] // MovementComponentMaker, returns a MovementComponent
		if !ok {
			panic(fmt.Sprintf("No registered maker for the component %s", key))
		}
		component := maker(data)
		e.ObjectPool.addComponent(component)
	}

	e.Entities = append(e.Entities, Entity{Index: e.lastActiveEntity, Name: data["UnitName"].(string)})

	e.lastActiveEntity++
}

//RemoveEntity WIP
func (e *EntityManager) RemoveEntity() {
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

//GetEntitiesData gets the data of all entities and packs them into []byte, used to send the clients necessary data to reconstruct the current state of the game
//TODO: add batching instead of sending all the data at once
func (e *EntityManager) GetEntitiesData() ([]byte, error) {
	size := e.lastActiveEntity
	entities := make([]EntityData, 0, size+1)

	for i := 0; i < size; i++ {
		//	fmt.Printf("I at %v am at position %v \n", i, e.PositionComponents[i].Position)
		entities = append(entities, EntityData{
			Index: i,
		})
	}

	data, err := json.Marshal(&entities)
	return data, err
}

//Used to sort actions by priority so we will save memory with CPU caching as the actions will be of the same type
type sortByActionPriority []Action

func (a sortByActionPriority) Len() int           { return len(a) }
func (a sortByActionPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByActionPriority) Less(i, j int) bool { return a[i].GetPriority() < a[j].GetPriority() }
