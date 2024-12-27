package main

import (
	"fmt"
	"game/internal/core"
	"game/internal/game"
	"game/internal/game/states"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	kernel := core.NewGameKernel()

	gameInstance := game.NewGame(kernel)

	kernel.EventBus.Subscribe("GameOver", func(data interface{}) {
		fmt.Println("Game over", data)

		gameInstance.SetState(states.MenuState)
	})

	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Survivor Game")

	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
