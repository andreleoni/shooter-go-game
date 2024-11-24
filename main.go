package main

import (
	"game/internal/core"
	"game/internal/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	kernel := core.NewGameKernel()

	gameInstance := game.NewGame(kernel)

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Survivor Game")
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
