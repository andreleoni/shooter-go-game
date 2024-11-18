package combat

import (
	"game/internal/core"
	"game/internal/helpers/collision"
	"game/internal/plugins"
	"game/internal/plugins/enemy"
	"game/internal/plugins/weapon"
	"game/internal/plugins/weapon/templates"

	"github.com/hajimehoshi/ebiten/v2"
)

type CombatPlugin struct {
	kernel      *core.GameKernel
	enemyPlugin *enemy.EnemyPlugin
}

func NewCombatPlugin(enemyPlugin *enemy.EnemyPlugin) *CombatPlugin {
	return &CombatPlugin{
		enemyPlugin: enemyPlugin,
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
	wp := cp.kernel.PluginManager.GetPlugin("WeaponSystem").(*weapon.WeaponPlugin)
	enemies := cp.enemyPlugin.GetEnemies()

	for _, weapon := range wp.GetWeapons() {
		for _, enemy := range enemies {
			for _, projectil := range weapon.Projectiles {
				if enemy.Active {
					if collision.Check(projectil.X, projectil.Y, 5, 10, enemy.X, enemy.Y, 20, 20) {
						projectil.Active = false

						enemy.Health -= projectil.Power

						if enemy.Health <= 0 {
							enemy.Active = false
						}
					}
				}
			}

			if weapon.Type == templates.ProtectionWeapon {
				if enemy.Active {
					playerPlugin := cp.kernel.PluginManager.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)
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
							enemy.DamageFlashTime = 0.1

							if enemy.Health <= 0 {
								enemy.Active = false
							}

						} else {
							enemy.LastProtectionDeltaTime += cp.kernel.DeltaTime
						}
					}
				}
			}
		}
	}

	return nil
}
