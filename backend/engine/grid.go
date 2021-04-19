package engine

type Grid interface {
	Update()
	GetEnemyProximityImap(tag int) *Imap
}
