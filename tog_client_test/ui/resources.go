package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	_ "image/png"
	"log"
)

var (
	P0_20Image *ebiten.Image
	P0_16Image *ebiten.Image
	P0_10Image *ebiten.Image

	P1_20Image *ebiten.Image
	P1_16Image *ebiten.Image
	P1_10Image *ebiten.Image
	BasicFont  font.Face

	StatsBarBg *ebiten.Image
)

func LoadResources() {
	loadFont()
	loadUnitSprites()

	var err error
	StatsBarBg, _, err = ebitenutil.NewImageFromFile("./resources/statsBarBg.png")
	check(err)
}

func loadUnitSprites()  {
	var err error
	P0_16Image, _, err = ebitenutil.NewImageFromFile("./resources/p0_16.png")
	P0_10Image, _, err = ebitenutil.NewImageFromFile("./resources/p0_10.png")
	P0_20Image, _, err = ebitenutil.NewImageFromFile("./resources/p0_20.png")
	check(err)

	P1_16Image, _, err = ebitenutil.NewImageFromFile("./resources/p1_16.png")
	P1_10Image, _, err = ebitenutil.NewImageFromFile("./resources/p1_10.png")
	P1_20Image, _, err = ebitenutil.NewImageFromFile("./resources/p1_20.png")
	check(err)
}

func loadFont(){
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	BasicFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    10,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}
func check(err error) {
	if err == nil {
		return
	}
	log.Fatalln(err)
}