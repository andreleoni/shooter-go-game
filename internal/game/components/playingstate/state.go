package playingstate

import (
	"fmt"
	"game/internal/core"
	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/chooseability"
	"game/internal/plugins/playing/combat"
	"game/internal/plugins/playing/enemy"
	"game/internal/plugins/playing/experience"
	"game/internal/plugins/playing/player"
	"game/internal/plugins/playing/stats"
	"game/internal/plugins/playing/weapon"
	"game/internal/plugins/playing/weapon/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

type State int

const (
	Playing State = iota
	Paused
	ChooseAbility
)

type ComponentPlayingState struct {
	kernel               *core.GameKernel
	pluginManagerByState map[State]*core.PluginManager
	state                State
}

func NewComponentPlayingState(kernel *core.GameKernel) *ComponentPlayingState {
	componentPlayingState := &ComponentPlayingState{kernel: kernel}

	pluginManagerByState := make(map[State]*core.PluginManager)
	pluginManagerByState[Playing] = core.NewPluginManager()
	pluginManagerByState[ChooseAbility] = core.NewPluginManager()

	kernel.EventBus.Subscribe("StartGame", func(data interface{}) {
		pluginManagerByState[Playing].UnregisterAll()

		playerPlugin := player.NewPlayerPlugin(pluginManagerByState[Playing])
		cameraPlugin := camera.NewCameraPlugin(playerPlugin)
		enemyPlugin := enemy.NewEnemyPlugin(playerPlugin, pluginManagerByState[Playing])
		combatPlugin := combat.NewCombatPlugin(enemyPlugin, pluginManagerByState[Playing])
		statsPlugin := stats.NewStatsPlugin(playerPlugin)
		weaponPlugin := weapon.NewWeaponPlugin(pluginManagerByState[Playing])
		experiencePlugin := experience.NewExperiencePlugin(pluginManagerByState[Playing])

		pluginManagerByState[Playing].Register(weaponPlugin, 0)
		pluginManagerByState[Playing].Register(playerPlugin, 1)
		pluginManagerByState[Playing].Register(experiencePlugin, 2)
		pluginManagerByState[Playing].Register(enemyPlugin, 3)
		pluginManagerByState[Playing].Register(combatPlugin, 4)
		pluginManagerByState[Playing].Register(cameraPlugin, 6)
		pluginManagerByState[Playing].Register(statsPlugin, 7)

		playerPlugin.Init(kernel)
		experiencePlugin.Init(kernel)
		enemyPlugin.Init(kernel)
		combatPlugin.Init(kernel)
		cameraPlugin.Init(kernel)
		statsPlugin.Init(kernel)
		weaponPlugin.Init(kernel)

		kernel.EventBus.Subscribe("NewAbility", func(weapon interface{}) {
			weaponType := weapon.(entities.WeaponType)
			weaponPlugin.AddWeapon(&entities.Weapon{Power: 20, Type: weaponType})

			componentPlayingState.SetState(Playing)
		})
	})

	kernel.EventBus.Subscribe("ChoosingAbility", func(data interface{}) {
		fmt.Println("Choosing ability", data)

		pluginManagerByState[ChooseAbility].UnregisterAll()

		chooseabilityPlugin := chooseability.NewChooseAbilityPlugin(pluginManagerByState[ChooseAbility])

		pluginManagerByState[ChooseAbility].Register(chooseabilityPlugin, 0)

		chooseabilityPlugin.Init(kernel)

		componentPlayingState.SetState(ChooseAbility)
	})

	componentPlayingState.pluginManagerByState = pluginManagerByState

	return componentPlayingState
}

func (cps *ComponentPlayingState) Draw(screen *ebiten.Image) {
	cps.pluginManagerByState[cps.state].DrawAll(screen)
}

func (cps *ComponentPlayingState) PluginManager() *core.PluginManager {
	return cps.pluginManagerByState[cps.state]
}

func (cps *ComponentPlayingState) SetState(s State) {
	cps.state = s
}
