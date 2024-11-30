package combat

import (
	"game/internal/core"
	"game/internal/helpers/collision"
	"game/internal/plugins"
	"game/internal/plugins/playing/enemy"
	"game/internal/plugins/playing/experience"
	"game/internal/plugins/playing/weapon"
	"game/internal/plugins/playing/weapon/templates"

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

	enemies := cp.enemyPlugin.GetEnemies()

	enemyGotDamaged := false

	for _, weapon := range wp.GetWeapons() {
		for _, enemy := range enemies {
			for _, projectil := range weapon.Projectiles {
				if enemy.Active {
					if collision.Check(projectil.X, projectil.Y, 5, 10, enemy.X, enemy.Y, enemy.Width, enemy.Height) {
						projectil.Active = false

						enemy.Health -= projectil.Power

						if enemy.Health <= 0 {
							enemy.Active = false
							enemyGotDamaged = true

							ep.DropCrystal(enemy.X, enemy.Y)
						}
					}
				}
			}

			if weapon.Type == templates.ProtectionWeapon {
				if enemy.Active {
					playerPlugin := cp.plugins.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)
					playerX, playerY := playerPlugin.GetPosition()

					if collision.CheckCircle(
						playerX,
						playerY,
						50,
						enemy.X,
						enemy.Y,
						enemy.Width,
						enemy.Height) {

						if enemy.LastProtectionDeltaTime >= 0.5 {
							enemy.Health -= weapon.Power
							enemy.LastProtectionDeltaTime = 0

							if enemy.Health <= 0 {
								enemy.Active = false

								ep.DropCrystal(enemy.X, enemy.Y)
							} else {
								enemyGotDamaged = true
							}

						} else {
							enemy.LastProtectionDeltaTime += cp.kernel.DeltaTime
						}
					}
				}
			}

			if enemyGotDamaged {
				enemy.DamageFlashTime = 0.1
			}
		}
	}

	return nil
}
