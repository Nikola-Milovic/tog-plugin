package tests

import (
	"testing"

	"github.com/Nikola-Milovic/tog-plugin/game"
)

func TestHandlers(t *testing.T) {
	world := game.CreateWorld()
	println(world)
}
