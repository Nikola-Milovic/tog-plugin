package engine

// AI represents interface which every AI has to inherit, its job is to produce the best action that unit should execute
type AI interface {
	PerformAI(index int)
}
