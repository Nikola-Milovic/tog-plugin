package ai

import (
	"fmt"
	"github.com/Nikola-Milovic/tog-plugin/constants"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/math"
)

type GenericAI struct {
	World *game.World
}

func (ai GenericAI) PerformAI(index int) {
	w := ai.World
	//g := ai.World.Grid
	entities := w.EntityManager.GetEntities()
	atkComp := w.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	posComp := w.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	movComp := w.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)

	playerTag := entities[index].PlayerTag
	attackRange := atkComp.Range
	engagingDistance := movComp.MovementSpeed*20 + attackRange

	fmt.Printf("I %d am at state %s\n", index, entities[index].State)

	switch entities[index].State {
	case constants.StateWalking:
		{
			closestEnemyIndex := -1
			closestEnemyDistance := float32(10000)
			for eIndx, p := range w.ObjectPool.Components["PositionComponent"] {
				if entities[eIndx].PlayerTag == playerTag {
					continue
				}
				pos := p.(components.PositionComponent)
				dist := math.GetDistanceIncludingDiagonalVectors(pos.Position, posComp.Position) - posComp.BoundingBox.X/2 - pos.BoundingBox.X/2
				if dist <= attackRange {
					attackTarget(index, entities[eIndx].ID, entities[index].ID, w)
					return
				}

				if dist < closestEnemyDistance {
					closestEnemyDistance = dist
					closestEnemyIndex = eIndx
				}
			}
			if closestEnemyIndex == -1 {
				return
			}

			if closestEnemyDistance <= engagingDistance {
				moveTowardsTarget(index, entities[closestEnemyIndex].ID, w, true)
			} else {
				moveTowardsTarget(index, entities[closestEnemyIndex].ID, w, false)
			}
		}
	case constants.StateEngaging:
		{

			targetIndex := w.EntityManager.GetIndexMap()[movComp.TargetID]
			targetPos := w.ObjectPool.Components["PositionComponent"][targetIndex].(components.PositionComponent)

			dist := math.GetDistanceIncludingDiagonalVectors(targetPos.Position, posComp.Position) - posComp.BoundingBox.X/2 - targetPos.BoundingBox.X/2
			if  dist > engagingDistance {
				entities[index].State = constants.StateWalking
				return
			} else if dist <= attackRange {
				attackTarget(index, movComp.TargetID, entities[index].ID, w)
			}
		}
	case constants.StateAttacking:
		{
			return
		}
	}

}

//if atkComp.IsAttacking {
//return
//}
//
//attackRange := atkComp.Range*10
//
//// 1) Create working map size of interest map
//proximity := 2 * int(movComp.Velocity.Magnitute()) * 5
//if proximity == 0 {
//proximity = 6 * int(movComp.MovementSpeed)
//}
//
//sizeX := proximity
//sizeY := proximity
//
//wm := g.GetWorkingMap(sizeX, sizeY)
//
//playerTag := entities[index].PlayerTag
//
//// 2) Add enemies with -1 mag and allies with 1.1 mag, add positions with -3 mag
//x, y := grid.GlobalCordToTiled(posComp.Position)
//engine.AddMaps(g.GetEnemyProximityImap(playerTag), wm, x, y, 1)
//
//targetX, targetY, value := wm.GetHighestCell()
//
//if value == 0.0 { // no enemies
//targetX, targetY, _ = g.GetEnemyProximityImap(playerTag).GetHighestCell()
//translatedX, translatedY := translateCoordsOutsideofMapIntoMap(x, y, targetX, targetY, sizeX)
//engine.AddMaps(grid.GetProximityTemplate(float32(8)).Imap, wm, translatedX, translatedY, 2)
//engine.AddMaps(g.GetProximityImaps()[playerTag], wm, x, y, -1)
//wm.NormalizeAndInvert()
//engine.AddMaps(g.GetInterestTemplate(sizeX), wm, x, y, 1)
//targetX, targetY, _ = wm.GetLowestValue()
//moveTowardsTarget(index, targetX, targetY, w)
//return
//} else {
//for eIndx, p := range w.ObjectPool.Components["PositionComponent"] {
//if entities[eIndx].PlayerTag == playerTag {
//continue
//}
//pos := p.(components.PositionComponent)
//dist := engine.GetDistanceIncludingDiagonalVectors(pos.Position, posComp.Position)
//if dist-posComp.BoundingBox.X/2-pos.BoundingBox.X/2 <= attackRange {
//attackTarget(index, entities[eIndx].ID, entities[index].ID, w)
//return
//}
//}
//
//engine.AddMaps(g.GetProximityImaps()[playerTag], wm, x, y, -1)
//engine.AddMaps(g.GetProximityImaps()[playerTag], wm, x, y, -0.5)
//wm.NormalizeAndInvert()
//engine.AddMaps(g.GetInterestTemplate(sizeX), wm, x, y, 1)
//targetX, targetY, _ = wm.GetLowestValue()
//
//moveTowardsTarget(index, targetX, targetY, w)
//return
//}
