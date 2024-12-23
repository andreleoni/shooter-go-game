package weapon

import (
	"game/internal/core"
	"game/internal/plugins"
	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/player"

	entitiesweapon "game/internal/plugins/playing/weapon/entities/weapons"

	"github.com/hajimehoshi/ebiten/v2"
)

type WeaponPlugin struct {
	kernel  *core.GameKernel
	plugins *core.PluginManager

	weapons []entitiesweapon.Weapon
}

func NewWeaponPlugin(plugins *core.PluginManager) *WeaponPlugin {
	return &WeaponPlugin{
		plugins: plugins,
	}
}

func (wp *WeaponPlugin) ID() string {
	return "WeaponSystem"
}

func (wp *WeaponPlugin) Init(
	kernel *core.GameKernel,
) error {

	wp.kernel = kernel

	return nil
}

func (wp *WeaponPlugin) AddWeapon(weapon entitiesweapon.Weapon) {
	wp.weapons = append(wp.weapons, weapon)
}

func (wp *WeaponPlugin) GetWeapons() []entitiesweapon.Weapon {
	return wp.weapons
}

func (wp *WeaponPlugin) Update() error {
	playerPlugin := wp.plugins.GetPlugin("PlayerSystem").(*player.PlayerPlugin)
	playerX, playerY := playerPlugin.GetPosition()

	cameraPlugin := wp.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for _, weapon := range wp.weapons {
		wui := entitiesweapon.WeaponUpdateInput{
			DeltaTime: wp.kernel.DeltaTime,
			PlayerX:   playerX,
			PlayerY:   playerY,
			CameraX:   cameraX,
			CameraY:   cameraY,
		}

		weapon.Update(wui)
	}

	return nil
}

func (wp *WeaponPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := wp.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	playerPlugin := wp.plugins.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)
	playerX, playerY := playerPlugin.GetPosition()

	wdi := entitiesweapon.WeaponDrawInput{
		CameraX: cameraX,
		CameraY: cameraY,
		PlayerX: playerX,
		PlayerY: playerY,
	}

	for _, weapon := range wp.weapons {
		weapon.Draw(screen, wdi)
	}
}
