package combat

import (
	"game/internal/core"
	"game/internal/plugins"
	"game/internal/plugins/playing/ability"
	"game/internal/plugins/playing/enemy"
	"game/internal/plugins/playing/experience"

	entitiesabilities "game/internal/plugins/playing/ability/entities/abilities"

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

	for _, a := range wp.GetAcquiredAbilities() {
		for _, enemy := range enemies {
			combatOutput := a.Combat(entitiesabilities.CombatInput{
				DeltaTime:    cp.kernel.DeltaTime,
				Enemy:        enemy,
				PlayerPlugin: pp,
				EnemyPlugin:  cp.enemyPlugin,
			})

			if combatOutput.EnemyGotDamaged {
				cp.enemyPlugin.ApplyDamage(
					enemy,
					combatOutput.Damage,
					combatOutput.CriticalDamage)

				enemy.DamageFlashTime = 0.1
			}

			if enemy.Active {
				if enemy.Health <= 0 {
					enemy.Active = false

					cp.enemyPlugin.AddDeathEnemies(enemy)
					ep.DropCrystal(enemy.X, enemy.Y)
				}
			}
		}
	}

	return nil
}
