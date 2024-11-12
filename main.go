package main

import (
	"game/internal/core"
	"game/internal/plugins/bullet"
	"game/internal/plugins/collision"
	"game/internal/plugins/enemy"
	"game/internal/plugins/obstacle"
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
	bulletPlugin := bullet.NewBulletPlugin()
	enemyPlugin := enemy.NewEnemyPlugin(playerPlugin)
	collisionPlugin := collision.NewCollisionPlugin(bulletPlugin, enemyPlugin)
	obstaclePlugin := obstacle.NewObstaclePlugin()

	kernel.PluginManager.Register(playerPlugin)
	kernel.PluginManager.Register(bulletPlugin)
	kernel.PluginManager.Register(enemyPlugin)
	kernel.PluginManager.Register(collisionPlugin)
	kernel.PluginManager.Register(obstaclePlugin)

	playerPlugin.Init(kernel)
	bulletPlugin.Init(kernel)
	enemyPlugin.Init(kernel)
	collisionPlugin.Init(kernel)
	obstaclePlugin.Init(kernel)

	game := &Game{kernel: kernel}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Survivor Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
