package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/engine"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/grid"
)

type GenericAI struct {
	World *game.World
}

func (ai GenericAI) PerformAI(index int) {
	w := ai.World
	g := ai.World.Grid
	entities := w.EntityManager.GetEntities()
	//atkComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	posComp := w.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	movComp := w.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)

	// 1) Create working map size of interest map
	proximity := 2 * int(movComp.Velocity.Magnitute()) * 5
	if proximity == 0 {
		proximity = 6 * int(movComp.MovementSpeed)
	}

	sizeX := grid.MapWidth/ grid.TileSize -1
	sizeY := grid.MapHeight/ grid.TileSize -1

	wm := g.GetWorkingMap(sizeX, sizeY)

	// 2) Add enemies with -1 mag and allies with 1.1 mag, add positions with -3 mag
	x, y := grid.GlobalCordToTiled(posComp.Position)
	engine.AddMaps(g.GetEnemyProximityImap(entities[index].PlayerTag), wm, x, y, 1)
	engine.AddMaps(g.GetProximityImaps()[entities[index].PlayerTag], wm, x, y, 0.5)
	engine.AddMaps(grid.GetProximityTemplate(movComp.MovementSpeed).Imap, wm, x,y, -1)
	wm.NormalizeAndInvert()
	engine.AddMaps(g.GetInterestTemplate(proximity), wm, x, y, 1)

	targetX, targetY := wm.GetLowestValue()
	moveTowardsPoint(index, targetX, targetY, w)

	//engine.PrintImapToFile(wm, "Workingmap", true)

	// 3) If concenration of enemies > x go towards it

	// 4) else go towards the place with the lowest ally count

	//	adjustedAtkRange := atkComp.Range*16 + posComp.BoundingBox.DivideScalar(2).X + 2
}
