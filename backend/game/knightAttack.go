package game

type BasicAttack interface {
	BasicAttack()
}

type KnightBasicAttack struct {
}

func (ba *KnightBasicAttack) BasicAttack() {}
