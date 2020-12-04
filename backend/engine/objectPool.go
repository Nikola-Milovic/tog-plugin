package engine

type ObjectPool struct {
	Components  map[string][]Component
	MaxSize     int
	currentSize int
}

func CreateObjectPool(maxSize int) *ObjectPool {
	op := ObjectPool{MaxSize: maxSize,
		Components: make(map[string][]Component, 10),
	}
	return &op
}

func (op *ObjectPool) addComponent(comp Component) {
	_, ok := op.Components[comp.ComponentName()]
	if !ok {
		op.Components[comp.ComponentName()] = make([]Component, 0, op.MaxSize)
	}

	op.Components[comp.ComponentName()] = append(op.Components[comp.ComponentName()], comp)
}
