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

	abilitiesrepository "game/internal/plugins/playing/ability/repository"

	"github.com/hajimehoshi/ebiten/v2"
)

type AbilityPlugin struct {
	kernel  *core.GameKernel
	plugins *core.PluginManager

	availableAbilities *abilitiesrepository.Ability
	acquiredAbilities  *abilitiesrepository.Ability
}

func NewAbilityPlugin(plugins *core.PluginManager) *AbilityPlugin {
	return &AbilityPlugin{
		plugins: plugins,
	}
}

func (wp *AbilityPlugin) ID() string {
	return "AbilitySystem"
}

func (ap *AbilityPlugin) Init(
	kernel *core.GameKernel,
) error {
	ap.availableAbilities = abilitiesrepository.NewAbility()
	ap.acquiredAbilities = abilitiesrepository.NewAbility()

	ap.kernel = kernel

	ap.AddAvailableAbility(abilitiesentitiesbasic.New())
	ap.AddAvailableAbility(abilitiesentitiesdagger.New())
	ap.AddAvailableAbility(abilitiesentitiesprotection.New())

	return nil
}

func (wp *AbilityPlugin) AcquireAbility(a abilitiesentities.Ability) {
	wp.acquiredAbilities.Add(a)
}

func (wp *AbilityPlugin) GetAcquiredAbilities() []abilitiesentities.Ability {
	return wp.acquiredAbilities.Get()
}

func (wp *AbilityPlugin) AddAvailableAbility(a abilitiesentities.Ability) {
	wp.availableAbilities.Add(a)
}

func (wp *AbilityPlugin) GetAvailableAbilitiesByName(a string) abilitiesentities.Ability {
	for _, wa := range wp.availableAbilities.Get() {
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

	for _, a := range wp.acquiredAbilities.Get() {
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

	for _, a := range wp.acquiredAbilities.Get() {
		a.Draw(screen, wdi)
	}
}
