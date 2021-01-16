package engine

type WorldI interface {
	GetObjectPool() *ObjectPool
	GetEntityManager() *EntityManager
	GetUnitDataMap() map[string]map[string]interface{}
	GetAbilityDataMap() map[string]map[string]interface{}
	GetEffectDataMap() map[string]map[string]interface{}
}
