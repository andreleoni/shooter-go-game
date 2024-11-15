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
	"time"

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

		kernel.PluginManager.Register(playerPlugin, 0)
		kernel.PluginManager.Register(bulletPlugin, 1)
		kernel.PluginManager.Register(enemyPlugin, 2)
		kernel.PluginManager.Register(combatPlugin, 3)
		kernel.PluginManager.Register(obstaclePlugin, 4)
		kernel.PluginManager.Register(cameraPlugin, 5)

		playerPlugin.Init(kernel)
		bulletPlugin.Init(kernel)
		enemyPlugin.Init(kernel)
		combatPlugin.Init(kernel)
		obstaclePlugin.Init(kernel)
		cameraPlugin.Init(kernel)

		time.Sleep(100 * time.Millisecond)
	})

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Survivor Game")
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
