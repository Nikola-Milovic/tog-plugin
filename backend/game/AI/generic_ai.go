package ai

import (
	"github.com/Nikola-Milovic/tog-plugin/game"
)

type GenericAI struct {
	World *game.World
}

func (ai GenericAI) PerformAI(index int) {
	//world := ai.World
	////g := ai.World.Grid
	////indexMap := world.EntityManager.GetIndexMap()
	////entities := world.EntityManager.GetEntities()
	////atkComp := world.ObjectPool.Components["AttackComponent"][index].(components.AttackComponent)
	//////posComp := world.ObjectPool.Components["PositionComponent"][index].(components.PositionComponent)
	////movComp := world.ObjectPool.Components["MovementComponent"][index].(components.MovementComponent)
	//
	//if entities[index].State != constants.StateThinking {
	//	return
	//}
	////
	////unitID := entities[index].ID
	////
	////playerTag := entities[index].PlayerTag
	////attackRange := atkComp.Range
	////engagingDistance := movComp.MovementSpeed*20 + attackRange
	//
	//switch entities[index].State {
	//case constants.StateWalking:
	//	{
	//
	//	}
	//case constants.StateEngaging:
	//	{
	//
	//	}
	//case constants.StateAttacking:
	//	{
	//
	//	}
	//case constants.StateThinking:
	//	{
	//	}
	//}
}