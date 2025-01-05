package playingstate

import (
	"fmt"
	"game/internal/core"
	"game/internal/plugins/playing/ability"
	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/chooseability"
	"game/internal/plugins/playing/combat"
	"game/internal/plugins/playing/enemy"
	"game/internal/plugins/playing/experience"
	"game/internal/plugins/playing/player"
	"game/internal/plugins/playing/scenario"
	"game/internal/plugins/playing/stats"

	entitiesability "game/internal/plugins/playing/ability/entities/abilities"
	playerentities "game/internal/plugins/playing/player/entities"

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

		character := data.(playerentities.Character)

		// Playing plugin
		playerPlugin := player.NewPlayerPlugin(pluginManagerByState[Playing], character)
		cameraPlugin := camera.NewCameraPlugin(playerPlugin)
		enemyPlugin := enemy.NewEnemyPlugin(playerPlugin, pluginManagerByState[Playing])
		combatPlugin := combat.NewCombatPlugin(enemyPlugin, pluginManagerByState[Playing])
		statsPlugin := stats.NewStatsPlugin(pluginManagerByState[Playing])
		abilityPlugin := ability.NewAbilityPlugin(pluginManagerByState[Playing])
		experiencePlugin := experience.NewExperiencePlugin(pluginManagerByState[Playing])
		scenarioPlugin := scenario.New(pluginManagerByState[Playing])

		pluginManagerByState[Playing].Register(scenarioPlugin, 1)
		pluginManagerByState[Playing].Register(abilityPlugin, 10)
		pluginManagerByState[Playing].Register(playerPlugin, 20)
		pluginManagerByState[Playing].Register(experiencePlugin, 30)
		pluginManagerByState[Playing].Register(enemyPlugin, 40)
		pluginManagerByState[Playing].Register(combatPlugin, 50)
		pluginManagerByState[Playing].Register(cameraPlugin, 60)
		pluginManagerByState[Playing].Register(statsPlugin, 70)

		playerPlugin.Init(kernel)
		experiencePlugin.Init(kernel)
		enemyPlugin.Init(kernel)
		combatPlugin.Init(kernel)
		cameraPlugin.Init(kernel)
		statsPlugin.Init(kernel)
		abilityPlugin.Init(kernel)
		scenarioPlugin.Init(kernel)

		// ChooseAbility plugins
		chooseabilityPlugin := chooseability.NewChooseAbilityPlugin(pluginManagerByState[ChooseAbility])

		pluginManagerByState[ChooseAbility].Register(chooseabilityPlugin, 0)

		chooseabilityPlugin.Init(kernel)

		// Subscribers
		kernel.EventBus.Subscribe("ChoosingAbility", func(data interface{}) {
			fmt.Println("ChoosingAbility")

			componentPlayingState.SetState(ChooseAbility)
		})

		kernel.EventBus.Subscribe("NewAbility", func(a interface{}) {
			ability := a.(entitiesability.Ability)
			ability.SetPluginManager(pluginManagerByState[Playing])
			abilityPlugin.AcquireAbility(ability)

			componentPlayingState.SetState(Playing)
		})

		kernel.EventBus.Publish(
			"NewAbility",
			abilityPlugin.GetAvailableAbilitiesByName(character.Ability))
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
