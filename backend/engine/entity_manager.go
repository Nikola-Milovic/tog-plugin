package engine

type EntityManagerI interface {
	Update()

	SetObjectPool(op *ObjectPool)
	SetEventManager(em *EventManager)
	SetComponentMaker(cm ComponentMaker)

	AddEntity(data NewEntityData, tag int, startOfMatch bool) (int, int)
	RemoveEntity(id int)
	GetIndexMap() map[int]int
	GetEntities() []Entity

	RegisterComponentMaker(componentName string, cm ComponentMakerFun)
	RegisterHandler(eventName string, h EventHandler)
	RegisterAIMaker(unitName string, aiMaker func() AI)
	RegisterSystem(sys System)
	RegisterTempSystem(sysName string, sysMaker func(w interface{}) TempSystem)

	AddTempSystem(sysName string, data map[string]interface{}, world WorldI)
	RemoveTempSystem(name string)

	StartMatch()

	//Testing
	GetSystems() []System
}
