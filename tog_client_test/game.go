package main

import (
	"encoding/json"
	"fmt"
	"github.com/Nikola-Milovic/client-test/ui"
	"github.com/Nikola-Milovic/tog-plugin/game"
	"github.com/Nikola-Milovic/tog-plugin/game/components"
	"github.com/Nikola-Milovic/tog-plugin/game/systems"
	"github.com/Nikola-Milovic/tog-plugin/math"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"image/color"
)

type Game struct {
	world          *game.World
	selectedUnitID int
}

//var P1Units = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[{\"x\":6,\"y\":10}],\"knight\":[{\"x\":9,\"y\":7}, {\"x\":9,\"y\":3}]}}")
//var p2Units = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":9,\"y\":10}]}}")

type playerData struct {
	Name  string                   `json:"name"`
	Units map[string][]math.Vector `json:"units"`
}

func (g *Game) init() {

	P1units := make(map[string][]math.Vector, 10)
	P1units["knight"] = []math.Vector{{10,3} , {10,4}, {10,2}, {11,5}, {10,5}, {7,3}, {9, 5},
	}
	//P1units["s_wolf"] = []math.Vector{{4,4}}
	//P1units["gob_spear"] = []math.Vector{{8,4}}
	//P1units["archer"] =  []math.Vector{{0,4}}
	P1Data := playerData{"Lemi", P1units}

	P1, err := json.Marshal(P1Data)
	check(err)

	p2Units := make(map[string][]math.Vector, 10)
	p2Units["knight"] = []math.Vector{{7,3}, {9, 5}, {8,10}, {10,4}, {10,2}}
	p2Data := playerData{"Lemi2", p2Units}

	p2, err := json.Marshal(p2Data)
	check(err)

	g.world = CreateWorld(P1, p2)
	g.world.StartMatch()
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

var tick = 0

func (g *Game) Update() error {
	tick++
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.getEntityAtPosition(x, y)
	}

	if tick%10 != 0 {
		return nil
	}
	go g.world.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{
		R: 190,
		G: 190,
		B: 190,
		A: 255,
	})

	posComps := g.world.ObjectPool.Components["PositionComponent"]
	movComps := g.world.ObjectPool.Components["MovementComponent"]
	//atkComps := g.world.ObjectPool.Components["AttackComponent"]
	for index, ent := range g.world.EntityManager.GetEntities() {
		posComp := posComps[index].(components.PositionComponent)
		movComp := movComps[index].(components.MovementComponent)
	//	atkComp := atkComps[index].(components.AttackComponent)

		drawUnit(screen, posComp, ent.PlayerTag, index, ent.State, g.selectedUnitID == ent.ID)

		//ebitenutil.DrawLine(screen, float64(posComp.Position.X), float64(posComp.Position.Y), float64(posComp.Position.X+movComp.Velocity.X*10),
		//	float64(posComp.Position.Y+movComp.Velocity.Y*10),
		//	colornames.Aqua)
		//
		//ebitenutil.DrawLine(screen, float64(posComp.Position.X), float64(posComp.Position.Y), float64(posComp.Position.X+movComp.Separation.X*30),
		//	float64(posComp.Position.Y+movComp.Separation.Y*30),
		//	colornames.Yellow)
		//
		//ebitenutil.DrawLine(screen, float64(posComp.Position.X), float64(posComp.Position.Y), float64(posComp.Position.X+movComp.Avoidance.X * 30),
		//	float64(posComp.Position.Y+movComp.Avoidance.Y * 30),
		//	colornames.Pink)

		ebitenutil.DrawLine(screen, float64(posComp.Position.X), float64(posComp.Position.Y), float64(posComp.Position.X+movComp.Seek.X * 30),
			float64(posComp.Position.Y+movComp.Seek.Y * 30),
			colornames.Limegreen)


	}
	g.selectedUnit(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (sw, sh int) {
	return screenWidth, screenHeight
}

func (g *Game) drawTangents(screen *ebiten.Image, index int) {

	posComps := g.world.ObjectPool.Components["PositionComponent"]
	movComps := g.world.ObjectPool.Components["MovementComponent"]

	ent := g.world.EntityManager.GetEntities()[index]
	posComp := posComps[index].(components.PositionComponent)
	movComp := movComps[index].(components.MovementComponent)

	if movComp.Velocity == math.Zero() {
		return
	}

	targetPos := movComp.Goal
	toTarget := posComp.Position.To(targetPos)
	distanceToTarget := posComp.Position.Distance(targetPos)
	//	futurePos := posComp.Position.Add(toTarget.Normalize().MultiplyScalar(64))
	//	radius := posComp.Radius + 80
	//square := math.Square(futurePos, radius)

	square := systems.GetSpatalSquareDebug(posComp.Position, toTarget, posComp.Radius, distanceToTarget)
	g.world.Buff = g.world.SpatialHash.Query(square, g.world.Buff[:0], ent.PlayerTag)

	ebitenutil.DrawRect(screen, float64(square.Center().X)-float64(square.W()/2), float64(square.Center().Y)-float64(square.W()/2), float64(square.W()), float64(square.H()), color.RGBA{225, 225, 225, 60})

	for _, id := range g.world.Buff {
		if id == ent.ID {
			continue
		}
		otherIndex := g.world.EntityManager.GetIndexMap()[id]

		otherPosComp := posComps[otherIndex].(components.PositionComponent)

		if distanceToTarget < otherPosComp.Position.Distance(posComp.Position) {
			continue
		}

		tanA, tanB, found := systems.GetTangents(otherPosComp.Position, otherPosComp.Radius+posComp.Radius+4, posComp.Position)
		if found {
			ebitenutil.DrawLine(screen, float64(posComp.Position.X), float64(posComp.Position.Y), float64(tanA.X),
				float64(tanA.Y),
				color.RGBA{R: 255, A: 255})

			ebitenutil.DrawLine(screen, float64(posComp.Position.X), float64(posComp.Position.Y), float64(tanB.X),
				float64(tanB.Y),
				color.RGBA{R: 255, A: 255})
		}
	}
}

func drawUnit(dst *ebiten.Image, posComp components.PositionComponent, tag int, index int, state string, selected bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(posComp.Position.X-posComp.Radius), float64(posComp.Position.Y-posComp.Radius))
	op.Filter = ebiten.FilterLinear
	if selected {
		hue := float64(30) * 2 * 3.14 / 128
		saturation := float64(30) / 128
		value := float64(30) / 128
		op.ColorM.ChangeHSV(hue, saturation, value)
	}

	if tag == 0 {
		switch posComp.Radius {
		case 10:
			dst.DrawImage(ui.P0_10Image, op)
		case 16:
			dst.DrawImage(ui.P0_16Image, op)
		case 20:
			dst.DrawImage(ui.P0_20Image, op)
		}

		text.Draw(dst, fmt.Sprintf("%d", index), ui.BasicFont, int(posComp.Position.X-posComp.Radius), int(posComp.Position.Y-posComp.Radius), colornames.Red)
		text.Draw(dst, state, ui.BasicFont, int(posComp.Position.X), int(posComp.Position.Y), colornames.Red)
	} else {
		switch posComp.Radius {
		case 10:
			dst.DrawImage(ui.P1_10Image, op)
		case 16:
			dst.DrawImage(ui.P1_16Image, op)
		case 20:
			dst.DrawImage(ui.P1_20Image, op)
		}
		text.Draw(dst, fmt.Sprintf("%d", index), ui.BasicFont, int(posComp.Position.X-posComp.Radius), int(posComp.Position.Y-posComp.Radius), colornames.Blue)
		text.Draw(dst, state, ui.BasicFont, int(posComp.Position.X), int(posComp.Position.Y), colornames.Red)
	}
}

func (g *Game) selectedUnit(dst *ebiten.Image) {
	if g.selectedUnitID != 0 {
		op := &ebiten.DrawImageOptions{}
		w, h := dst.Size()
		op.GeoM.Translate(float64(w/2)-270, float64(h-108))
		op.Filter = ebiten.FilterLinear
		dst.DrawImage(ui.StatsBarBg, op)

		//	posComps := g.world.ObjectPool.Components["PositionComponent"]
		movComps := g.world.ObjectPool.Components["MovementComponent"]
		atkComps := g.world.ObjectPool.Components["AttackComponent"]
		statsComps := g.world.ObjectPool.Components["StatsComponent"]

		index := g.world.EntityManager.GetIndexMap()[g.selectedUnitID]

		g.drawTangents(dst, index)

		ent := g.world.EntityManager.GetEntities()

		if !ent[index].Active {
			g.selectedUnitID = 0
			return
		}

		//posComp := posComps[index].(components.PositionComponent)
		movComp := movComps[index].(components.MovementComponent)
		atkComp := atkComps[index].(components.AttackComponent)
		statsComp := statsComps[index].(components.StatsComponent)

		t := fmt.Sprintf("ID: %d,    Index: %d,    State: %s \n\n  Unit: %s,     MS: %.2f,      HP: %d/%d \n\n  ATK: %d,       RA: %.2f,      AS: %d \n\n",
			ent[index].ID, ent[index].Index, ent[index].State, ent[index].UnitID,
			movComp.MovementSpeed, statsComp.Health, statsComp.MaxHealth,
			atkComp.Damage, atkComp.Range, atkComp.AttackSpeed)
		text.Draw(dst, t, ui.BasicFont, w/2-240, h-90, colornames.Black)

		av := fmt.Sprintf("Avoidance: %v", movComp.Avoidance)
		sep := fmt.Sprintf("Seperation: %v", movComp.Separation)
		veloc := fmt.Sprintf("Velocity: %v", movComp.Velocity)
		seek := fmt.Sprintf("Seek: %v", movComp.Seek)


		text.Draw(dst, av, ui.BasicFont, 20, 20, colornames.Pink)
		text.Draw(dst, sep, ui.BasicFont, 20, 30, colornames.Yellow)
		text.Draw(dst, veloc, ui.BasicFont, 20, 40, colornames.Aqua)
		text.Draw(dst, seek, ui.BasicFont, 20, 50, colornames.Limegreen)

		text.Draw(dst, "G", ui.BasicFont, int(movComp.Goal.X), int(movComp.Goal.Y), color.Black)
	}
}

func (g *Game) getEntityAtPosition(x int, y int) {
	entities := make([]int, 3)
	entities = g.world.SpatialHash.Query(math.Square(math.Vector{X: float32(x), Y: float32(y)}, 16), entities[:0], -1)
	if len(entities) > 0 {
		g.selectedUnitID = entities[0]
	}
}
