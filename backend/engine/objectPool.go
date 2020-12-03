package engine

type ObjectPool struct {
	Components  map[string][]Component
	MaxSize     int
	currentSize int
}

func CreateObjectPool(maxSize int) *ObjectPool {
	op := ObjectPool{MaxSize: maxSize}
	return &op
}
