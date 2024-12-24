package ability

import (
	"fmt"
	"game/internal/core"
	"game/internal/plugins"
	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/player"

	abilitiesentities "game/internal/plugins/playing/ability/entities/abilities"
	abilitiesentitiesbasic "game/internal/plugins/playing/ability/entities/abilities/basic"
	abilitiesentitiesdagger "game/internal/plugins/playing/ability/entities/abilities/dagger"
	abilitiesentitiesprotection "game/internal/plugins/playing/ability/entities/abilities/protection"

	"github.com/hajimehoshi/ebiten/v2"
)

type AbilityPlugin struct {
	kernel  *core.GameKernel
	plugins *core.PluginManager

	availableAbilities []abilitiesentities.Ability
	acquiredAbilities  []abilitiesentities.Ability
}

func NewAbilityPlugin(plugins *core.PluginManager) *AbilityPlugin {
	return &AbilityPlugin{
		plugins: plugins,
	}
}

func (wp *AbilityPlugin) ID() string {
	return "AbilitySystem"
}

func (wp *AbilityPlugin) Init(
	kernel *core.GameKernel,
) error {

	wp.kernel = kernel

	wp.AddAvailableAbility(abilitiesentitiesbasic.New())
	wp.AddAvailableAbility(abilitiesentitiesdagger.New())
	wp.AddAvailableAbility(abilitiesentitiesprotection.New())

	return nil
}

func (wp *AbilityPlugin) AcquireAbility(a abilitiesentities.Ability) {
	wp.acquiredAbilities = append(wp.acquiredAbilities, a)
}

func (wp *AbilityPlugin) GetAcquiredAbilities() []abilitiesentities.Ability {
	return wp.acquiredAbilities
}

func (wp *AbilityPlugin) AddAvailableAbility(a abilitiesentities.Ability) {
	wp.availableAbilities = append(wp.availableAbilities, a)
}

func (wp *AbilityPlugin) GetAvailableAbilitiesByName(a string) abilitiesentities.Ability {
	for _, wa := range wp.availableAbilities {
		fmt.Println("getting available abilities", wa.ID(), a)
		if wa.ID() == a {
			return wa
		}
	}

	return nil
}

func (wp *AbilityPlugin) Update() error {
	playerPlugin := wp.plugins.GetPlugin("PlayerSystem").(*player.PlayerPlugin)
	playerX, playerY := playerPlugin.GetPosition()

	cameraPlugin := wp.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for _, a := range wp.acquiredAbilities {
		wui := abilitiesentities.AbilityUpdateInput{
			DeltaTime: wp.kernel.DeltaTime,
			PlayerX:   playerX,
			PlayerY:   playerY,
			CameraX:   cameraX,
			CameraY:   cameraY,
		}

		a.Update(wui)
	}

	return nil
}

func (wp *AbilityPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := wp.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	playerPlugin := wp.plugins.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)
	playerX, playerY := playerPlugin.GetPosition()

	wdi := abilitiesentities.AbilityDrawInput{
		CameraX: cameraX,
		CameraY: cameraY,
		PlayerX: playerX,
		PlayerY: playerY,
	}

	for _, a := range wp.acquiredAbilities {
		a.Draw(screen, wdi)
	}
}
