package main

import (
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/bullet"
	"game/internal/plugins/camera"
	"game/internal/plugins/combat"
	"game/internal/plugins/enemy"
	"game/internal/plugins/menu"
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
	return constants.ScreenWidth, constants.ScreenHeight
}

func main() {
	kernel := core.NewGameKernel()

	menuPlugin := menu.NewMenuPlugin()

	playerPlugin := player.NewPlayerPlugin()
	cameraPlugin := camera.NewCameraPlugin(playerPlugin)
	enemyPlugin := enemy.NewEnemyPlugin(playerPlugin)
	bulletPlugin := bullet.NewBulletPlugin()
	combatPlugin := combat.NewCombatPlugin(bulletPlugin, enemyPlugin)
	obstaclePlugin := obstacle.NewObstaclePlugin()

	kernel.PluginManager.Register(menuPlugin)
	kernel.PluginManager.Register(playerPlugin)
	kernel.PluginManager.Register(bulletPlugin)
	kernel.PluginManager.Register(enemyPlugin)
	kernel.PluginManager.Register(combatPlugin)
	kernel.PluginManager.Register(obstaclePlugin)
	kernel.PluginManager.Register(cameraPlugin)

	menuPlugin.Init(kernel)
	playerPlugin.Init(kernel)
	bulletPlugin.Init(kernel)
	enemyPlugin.Init(kernel)
	combatPlugin.Init(kernel)
	obstaclePlugin.Init(kernel)
	cameraPlugin.Init(kernel)

	game := &Game{kernel: kernel}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Survivor Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
