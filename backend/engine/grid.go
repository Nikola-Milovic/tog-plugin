package engine

import "github.com/Nikola-Milovic/tog-plugin/math"

type Grid interface {
	Update()
	GetEnemyProximityImap(tag int) *Imap
	GetOccupationalMap() *Imap
	GetWorkingMap(width, height int) *Imap
	GetProximityImaps() []*Imap
	GetGoalMap() *Imap
	GetInterestTemplate(size int) *Imap
	IsPositionFree(pos math.Vector, bbox math.Vector) bool
}
