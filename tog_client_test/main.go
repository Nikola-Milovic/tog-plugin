package main

import (
	"github.com/Nikola-Milovic/client-test/ui"
	"github.com/Nikola-Milovic/tog-plugin/startup"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 512
)

func main() {
	ebiten.SetWindowSize(screenWidth+300, screenHeight+300)
	ebiten.SetWindowResizable(true)
	ebiten.SetScreenClearedEveryFrame(false)


	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func init() {
	startup.ResourcesPath = "../backend/resources"
	go startup.StartUp(true)
	ui.LoadResources()
}




