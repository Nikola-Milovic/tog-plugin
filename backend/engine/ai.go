package engine

type AI interface {
	CalculateAction(index int, e *EntityManager) Action
}
