package playingstate

import (
	"fmt"
	"game/internal/core"
	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/combat"
	"game/internal/plugins/playing/enemy"
	"game/internal/plugins/playing/player"
	"game/internal/plugins/playing/stats"
	"game/internal/plugins/playing/weapon"

	"github.com/hajimehoshi/ebiten/v2"
)

type ComponentPlayingState struct {
	kernel        *core.GameKernel
	pluginManager *core.PluginManager
}

func NewComponentPlayingState(kernel *core.GameKernel) *ComponentPlayingState {
	pluginManager := core.NewPluginManager()

	kernel.EventBus.Subscribe("StartGame", func(data interface{}) {
		pluginManager.UnregisterAll()

		fmt.Println("Game started", data)
		// Level design com apenas os plugins necessários para aquele level

		playerPlugin := player.NewPlayerPlugin()
		cameraPlugin := camera.NewCameraPlugin(playerPlugin)
		enemyPlugin := enemy.NewEnemyPlugin(playerPlugin)
		combatPlugin := combat.NewCombatPlugin(enemyPlugin)
		statsPlugin := stats.NewStatsPlugin(playerPlugin)
		weaponPlugin := weapon.NewWeaponPlugin()

		pluginManager.Register(weaponPlugin, 0)
		pluginManager.Register(playerPlugin, 1)
		pluginManager.Register(enemyPlugin, 3)
		pluginManager.Register(combatPlugin, 4)
		pluginManager.Register(cameraPlugin, 6)
		pluginManager.Register(statsPlugin, 7)

		playerPlugin.Init(kernel)
		enemyPlugin.Init(kernel)
		combatPlugin.Init(kernel)
		cameraPlugin.Init(kernel)
		statsPlugin.Init(kernel)
		weaponPlugin.Init(kernel)
	})

	return &ComponentPlayingState{kernel: kernel}
}

func (cps *ComponentPlayingState) Update() {
	cps.pluginManager.UpdateAll()
}

func (cps *ComponentPlayingState) Draw(screen *ebiten.Image) {
	cps.pluginManager.DrawAll(screen)
}
