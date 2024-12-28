package combat

import (
	"game/internal/core"
	"game/internal/helpers/collision"
	"game/internal/plugins"
	"game/internal/plugins/playing/ability"
	"game/internal/plugins/playing/enemy"
	"game/internal/plugins/playing/experience"

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
	wp := cp.plugins.GetPlugin("AbilitySystem").(*ability.AbilityPlugin)
	ep := cp.plugins.GetPlugin("ExperienceSystem").(*experience.ExperiencePlugin)
	pp := cp.plugins.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)

	enemies := cp.enemyPlugin.GetEnemies()

	playerX, playerY := pp.GetPosition()

	for _, a := range wp.GetAcquiredAbilities() {
		for _, enemy := range enemies {
			enemyGotDamaged := false
			enemykilled := false

			if a.DamageType() == "projectil" {
				for _, projectil := range a.ActiveProjectiles() {
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

							damage, critical := pp.CalculateDamage(projectil.Power)

							cp.enemyPlugin.ApplyDamage(enemy, damage, critical)

							if enemy.Health <= 0 {
								enemy.Active = false
								enemykilled = true
								cp.enemyPlugin.AddDeathEnemies(enemy)

							} else {
								enemyGotDamaged = true
							}
						}
					}
				}
			}

			if a.DamageType() == "area" {
				if enemy.Active {
					abilityID := a.ID()

					if collision.CheckCircle(
						playerX,
						playerY,
						200, //implementar interface com radius
						enemy.X,
						enemy.Y,
						enemy.Width,
						enemy.Height) {

						lastAreaDamageDeltaTime, exists := enemy.LastAreaDamageDeltaTimeByAbility[abilityID]
						if !exists {
							lastAreaDamageDeltaTime = 0
						}

						if lastAreaDamageDeltaTime >= a.AttackSpeed() {
							damage, critical := pp.CalculateDamage(a.GetPower())
							cp.enemyPlugin.ApplyDamage(enemy, damage, critical)

							lastAreaDamageDeltaTime = 0

							if enemy.Health <= 0 {
								enemy.Active = false
								enemykilled = true
								cp.enemyPlugin.AddDeathEnemies(enemy)

							} else {
								enemyGotDamaged = true
							}

						} else {
							lastAreaDamageDeltaTime += cp.kernel.DeltaTime
						}

						enemy.LastAreaDamageDeltaTimeByAbility[abilityID] = lastAreaDamageDeltaTime
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
