package engine

type WorldI interface {
	GetObjectPool() *ObjectPool
	GetEntityManager() EntityManagerI
	GetUnitDataMap() map[string]map[string]interface{}
	GetAbilityDataMap() map[string]map[string]interface{}
	GetEffectDataMap() map[string]map[string]interface{}
}
