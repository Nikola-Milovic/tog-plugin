package engine

type ObjectPool struct {
	Components map[string][]Component
	//UniqueComponents map[string]map[string]Component
	AI          map[string]AI
	MaxSize     int
	currentSize int
}

//CreateObjectPool initializes a an Object pool with the maxSize specified
func CreateObjectPool(maxSize int) *ObjectPool {
	op := ObjectPool{MaxSize: maxSize,
		Components: make(map[string][]Component, maxSize),
		AI:         make(map[string]AI, maxSize),
		//UniqueComponents: make(map[string]map[string]Component, 2), //Components that not all entityPositions have, EG, ability component
	}
	return &op
}

func (op *ObjectPool) AddComponent(comp Component) {
	_, ok := op.Components[comp.ComponentName()]
	if !ok {
		op.Components[comp.ComponentName()] = make([]Component, 0, op.MaxSize)
	}

	op.Components[comp.ComponentName()] = append(op.Components[comp.ComponentName()], comp)
}

func (op *ObjectPool) RemoveAt(index int, id int) {
	for i, components := range op.Components {
		if index >= len(components) {
			continue
		}
		components[index] = components[len(components)-1]
		components = components[:len(components)-1]
		op.Components[i] = components
	}
}
