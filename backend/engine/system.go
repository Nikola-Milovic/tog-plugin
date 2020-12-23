package engine

//System represents a classical approach to data driven ECS, where systems run once in an update and do operations on data they are interested in
//we use a system to check for continuous effects, lifetimes, buffs/ debuffs/ durations/ death...
type System interface {
	Update()
}
