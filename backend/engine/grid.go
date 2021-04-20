package engine

type Grid interface {
	Update()
	GetEnemyProximityImap(tag int) *Imap
	GetOccupationalMap() *Imap
	GetWorkingMap(width, height int) *Imap
	GetProximityImaps() []*Imap
	GetInterestTemplate(size int) *Imap
}
