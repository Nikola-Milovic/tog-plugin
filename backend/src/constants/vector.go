package constants

type V2 struct {
	X, Y float32
}

func (a V2) Vector2() V2 {
	return V2{a.X, a.Y}
}
