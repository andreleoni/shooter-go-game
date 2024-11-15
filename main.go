package main

import (
	"fmt"
	"game/internal/core"
	"game/internal/core/game"
	"game/internal/plugins/bullet"
	"game/internal/plugins/camera"
	"game/internal/plugins/combat"
	"game/internal/plugins/enemy"
	"game/internal/plugins/obstacle"
	"game/internal/plugins/player"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	kernel := core.NewGameKernel()

	gameInstance := game.NewGame(kernel)

	kernel.EventBus.Subscribe("StartGame", func(data interface{}) {
		fmt.Println("Game started", data)
		// Level design com apenas os plugins necessários para aquele level

		playerPlugin := player.NewPlayerPlugin()
		cameraPlugin := camera.NewCameraPlugin(playerPlugin)
		enemyPlugin := enemy.NewEnemyPlugin(playerPlugin)
		bulletPlugin := bullet.NewBulletPlugin()
		combatPlugin := combat.NewCombatPlugin(bulletPlugin, enemyPlugin)
		obstaclePlugin := obstacle.NewObstaclePlugin()

		kernel.PluginManager.Register(playerPlugin)
		kernel.PluginManager.Register(bulletPlugin)
		kernel.PluginManager.Register(enemyPlugin)
		kernel.PluginManager.Register(combatPlugin)
		kernel.PluginManager.Register(obstaclePlugin)
		kernel.PluginManager.Register(cameraPlugin)

		playerPlugin.Init(kernel)
		bulletPlugin.Init(kernel)
		enemyPlugin.Init(kernel)
		combatPlugin.Init(kernel)
		obstaclePlugin.Init(kernel)
		cameraPlugin.Init(kernel)
	})

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Survivor Game")
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
