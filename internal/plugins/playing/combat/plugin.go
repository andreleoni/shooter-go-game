package combat

import (
	"game/internal/core"
	"game/internal/helpers/collision"
	"game/internal/plugins"
	"game/internal/plugins/playing/enemy"
	"game/internal/plugins/playing/experience"
	"game/internal/plugins/playing/weapon"

	"github.com/hajimehoshi/ebiten/v2"
)

type CombatPlugin struct {
	kernel      *core.GameKernel
	plugins     *core.PluginManager
	enemyPlugin *enemy.EnemyPlugin
}

func NewCombatPlugin(
	enemyPlugin *enemy.EnemyPlugin,
	plugins *core.PluginManager) *CombatPlugin {

	return &CombatPlugin{
		enemyPlugin: enemyPlugin,
		plugins:     plugins,
	}
}

func (cp *CombatPlugin) ID() string {
	return "CombatSystem"
}

func (cp *CombatPlugin) Init(kernel *core.GameKernel) error {
	cp.kernel = kernel
	return nil
}

func (cp *CombatPlugin) Draw(*ebiten.Image) {
	return
}

func (cp *CombatPlugin) Update() error {
	wp := cp.plugins.GetPlugin("WeaponSystem").(*weapon.WeaponPlugin)
	ep := cp.plugins.GetPlugin("ExperienceSystem").(*experience.ExperiencePlugin)
	pp := cp.plugins.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)

	enemies := cp.enemyPlugin.GetEnemies()

	playerX, playerY := pp.GetPosition()

	for _, weapon := range wp.GetWeapons() {
		weapon.AutoShot(cp.kernel.DeltaTime, playerX, playerY)

		for _, enemy := range enemies {
			enemyGotDamaged := false
			enemykilled := false

			if weapon.DamageType() == "projectil" {
				for _, projectil := range weapon.ActiveProjectiles() {
					if enemy.Active && projectil.Active {
						if collision.Check(
							projectil.X,
							projectil.Y,
							projectil.Width,
							projectil.Height,
							enemy.X,
							enemy.Y,
							enemy.Width,
							enemy.Height) {

							projectil.Active = false

							enemy.Health -= projectil.Power

							if enemy.Health <= 0 {
								enemy.Active = false
								enemykilled = true

							} else {
								enemyGotDamaged = true
							}
						}
					}
				}
			}

			if weapon.DamageType() == "area" {
				if enemy.Active {
					weaponID := weapon.ID()

					if collision.CheckCircle(
						playerX,
						playerY,
						200, //implementar interface com radius
						enemy.X,
						enemy.Y,
						enemy.Width,
						enemy.Height) {

						lastAreaDamageDeltaTime, exists := enemy.LastAreaDamageDeltaTimeByWeapon[weaponID]
						if !exists {
							lastAreaDamageDeltaTime = 0
						}

						if lastAreaDamageDeltaTime >= weapon.AttackSpeed() {
							enemy.Health -= weapon.GetPower()
							lastAreaDamageDeltaTime = 0

							if enemy.Health <= 0 {
								enemy.Active = false
								enemykilled = true

							} else {
								enemyGotDamaged = true
							}

						} else {
							lastAreaDamageDeltaTime += cp.kernel.DeltaTime
						}

						enemy.LastAreaDamageDeltaTimeByWeapon[weaponID] = lastAreaDamageDeltaTime
					}
				}
			}

			if enemykilled {
				ep.DropCrystal(enemy.X, enemy.Y)
			}

			if enemyGotDamaged {
				enemy.DamageFlashTime = 0.1
			}
		}
	}

	return nil
}
