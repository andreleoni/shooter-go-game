package combat

import (
	"game/internal/core"
	"game/internal/helpers/collision"
	"game/internal/plugins"
	"game/internal/plugins/playing/ability"
	"game/internal/plugins/playing/camera"
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
}

func (cp *CombatPlugin) Update() error {
	wp := cp.plugins.GetPlugin("AbilitySystem").(*ability.AbilityPlugin)
	ep := cp.plugins.GetPlugin("ExperienceSystem").(*experience.ExperiencePlugin)
	pp := cp.plugins.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)
	cameraPlugin := cp.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)

	enemies := cp.enemyPlugin.GetEnemies()
	playerX, playerY := pp.GetPosition()
	playerWidth, playerHeight := pp.GetSize()
	cameraX, cameraY := cameraPlugin.GetPosition()

	for _, a := range wp.GetAcquiredAbilities() {
		for _, enemy := range enemies {
			combatOutput := a.Combat(entitiesabilities.CombatInput{
				DeltaTime:    cp.kernel.DeltaTime,
				Enemy:        enemy,
				PlayerPlugin: pp,
				EnemyPlugin:  cp.enemyPlugin,
				CameraX:      cameraX,
				CameraY:      cameraY,
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
					ep.DropCrystal(
						enemy.X+(enemy.Width/2),
						enemy.Y+(enemy.Height/2))
				}
			}
		}

		// Update projectiles
		for _, p := range cp.enemyPlugin.GetGlobalProjectiles() {
			if p.Active {
				p.X += p.DirectionX * p.Speed * cp.kernel.DeltaTime
				p.Y += p.DirectionY * p.Speed * cp.kernel.DeltaTime

				// Check collision with player
				playerCollision := collision.Check(
					p.X, p.Y,
					p.Width, p.Height,
					(playerX - playerWidth/2), (playerY - playerHeight/2),
					playerWidth, playerHeight)

				if playerCollision {
					pp.ApplyDamage(p.Power)
					p.Active = false
				}
			}
		}
	}

	return nil
}
