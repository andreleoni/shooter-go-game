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
	"game/internal/plugins/stats"
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
		statsPlugin := stats.NewStatsPlugin(playerPlugin)

		kernel.PluginManager.Register(playerPlugin, 0)
		kernel.PluginManager.Register(bulletPlugin, 1)
		kernel.PluginManager.Register(enemyPlugin, 2)
		kernel.PluginManager.Register(combatPlugin, 3)
		kernel.PluginManager.Register(obstaclePlugin, 4)
		kernel.PluginManager.Register(cameraPlugin, 5)
		kernel.PluginManager.Register(statsPlugin, 6)

		playerPlugin.Init(kernel)
		bulletPlugin.Init(kernel)
		enemyPlugin.Init(kernel)
		combatPlugin.Init(kernel)
		obstaclePlugin.Init(kernel)
		cameraPlugin.Init(kernel)
		statsPlugin.Init(kernel)
	})

	kernel.EventBus.Subscribe("GameOver", func(data interface{}) {
		fmt.Println("Game over", data)

		gameInstance.SetStateGameOver()
	})

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Survivor Game")
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
