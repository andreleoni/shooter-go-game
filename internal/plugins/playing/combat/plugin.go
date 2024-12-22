package combat

import (
	"fmt"
	"game/internal/core"
	"game/internal/helpers/collision"
	"game/internal/plugins"
	"game/internal/plugins/playing/enemy"
	"game/internal/plugins/playing/experience"
	"game/internal/plugins/playing/weapon"
	"time"

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

	x, y := pp.GetPosition()

	for _, weapon := range wp.GetWeapons() {
		weapon.AutoShot(cp.kernel.DeltaTime, x, y)

		for _, enemy := range enemies {
			enemyGotDamaged := false

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

							ep.DropCrystal(enemy.X, enemy.Y)
						} else {
							enemyGotDamaged = true
						}
					}
				}
			}

			//if weapon.Type == templates.ProtectionWeapon
			if false {
				playerPlugin := cp.plugins.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)

				if enemy.Active {
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
							enemy.Health -= weapon.GetPower()
							enemy.LastProtectionDeltaTime = 0
							enemy.DamageFlashTime = 0.1

							if enemy.Health <= 0 {
								enemy.Active = false
								fmt.Println("Enemy killed by protection weapon", time.Now().Unix())

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
