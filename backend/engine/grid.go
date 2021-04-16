package engine

type Grid interface {
	Update()
	GetDesiredDirectionAt(pos Vector, tag int) Vector
	IsPositionFree(index int, positionToCheck Vector, boundingBox Vector) int
}
