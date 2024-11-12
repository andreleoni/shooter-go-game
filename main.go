// main.go
package main

import (
	"game/internal/core"
	"game/internal/plugins/player"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	kernel *core.GameKernel
}

func (g *Game) Update() error {
	return g.kernel.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.kernel.PluginManager.DrawAll(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func main() {
	kernel := core.NewGameKernel()

	// Register plugins
	playerPlugin := player.NewPlayerPlugin()
	kernel.PluginManager.Register(playerPlugin)
	playerPlugin.Init(kernel)

	game := &Game{kernel: kernel}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Vampire Survivors")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
