package basic

import (
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/helpers/collision"
	"game/internal/plugins"
	"game/internal/plugins/playing/ability/entities"
	abilityentities "game/internal/plugins/playing/ability/entities/abilities"
	"game/internal/plugins/playing/camera"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Basic struct {
	plugins     *core.PluginManager
	Projectiles []*entities.Projectile
	Power       float64

	// Shoot cooldown
	ShootTimer    float64
	ShootCooldown float64

	Level int
}

func New() *Basic {
	return &Basic{
		Power:         10,
		ShootCooldown: 1.0,
		Level:         1,
	}
}

func (b *Basic) SetPluginManager(plugins *core.PluginManager) {
	b.plugins = plugins
}

func (b *Basic) ID() string {
	return "Basic"
}

func (b *Basic) AutoShot(deltaTime, x, y float64) {
	b.ShootTimer += deltaTime

	if b.ShootTimer >= b.ShootCooldown {
		b.Shoot(x, y)
		b.ShootTimer = 0
	}
}

func (b *Basic) Shoot(x, y float64) {
	// Get enemy plugin to find closest enemy
	enemyPlugin := b.plugins.GetPlugin("EnemySystem").(plugins.EnemyPlugin)

	enemies := enemyPlugin.GetEnemies()

	if len(enemies) > 0 {
		// Find closest enemy
		closestEnemy := enemies[0]
		closestDist := math.MaxFloat64

		for _, enemy := range enemies {
			if !enemy.Active {
				continue
			}

			dx := (enemy.X + enemy.Width/2) - x
			dy := (enemy.Y + enemy.Height/2) - y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist < closestDist {
				closestDist = dist
				closestEnemy = enemy
			}
		}

		// Calcular direção
		dx := (closestEnemy.X + closestEnemy.Width/2) - x
		dy := (closestEnemy.Y + closestEnemy.Height/2) - y

		distance := math.Sqrt(dx*dx + dy*dy)

		dirX := dx / distance
		dirY := dy / distance

		// Create bullet targeting closest enemy
		projectile := &entities.Projectile{
			X:          x,
			Y:          y,
			Speed:      300,
			Active:     true,
			Power:      b.Power,
			TargetX:    closestEnemy.X + closestEnemy.Width/2,
			TargetY:    closestEnemy.Y + closestEnemy.Height/2,
			DirectionX: dirX,
			DirectionY: dirY,
			Height:     10,
			Width:      10,
		}

		b.Projectiles = append(b.Projectiles, projectile)
	}
}

func (b *Basic) Update(wui abilityentities.AbilityUpdateInput) {
	b.AutoShot(wui.DeltaTime, wui.PlayerX, wui.PlayerY)

	for _, projectile := range b.Projectiles {
		if projectile.Active {
			dx := projectile.DirectionX
			dy := projectile.DirectionY

			// Calculate distance
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance > 0 {
				// Normalize direction
				dx /= distance
				dy /= distance

				// Update position
				projectile.X += dx * projectile.Speed * wui.DeltaTime
				projectile.Y += dy * projectile.Speed * wui.DeltaTime
			}

			cameraPlugin := b.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
			cameraX, cameraY := cameraPlugin.GetPosition()
			// Check if projectile is too far from camera view
			screenX := projectile.X - cameraX
			screenY := projectile.Y - cameraY
			margin := float64(200) // Larger margin before deactivating

			if screenX < -margin ||
				screenX > constants.ScreenWidth+margin ||
				screenY < -margin ||
				screenY > constants.ScreenHeight+margin {

				projectile.Active = false
			}
		}
	}
}

func (b *Basic) Draw(screen *ebiten.Image, wdi abilityentities.AbilityDrawInput) {
	for _, projectile := range b.Projectiles {
		if projectile.Active {
			// Draw bullet relative to camera position
			screenX := projectile.X - wdi.CameraX
			screenY := projectile.Y - wdi.CameraY

			// Only draw if on screen
			if screenX >= -5 && screenX <= constants.ScreenWidth+5 &&
				screenY >= -5 && screenY <= constants.ScreenHeight+5 {

				vector.DrawFilledRect(
					screen,
					float32(screenX),
					float32(screenY),
					float32(projectile.Width),
					float32(projectile.Height),
					color.RGBA{200, 255, 0, 255},
					true)

				angle := math.Atan2(
					projectile.DirectionY, projectile.DirectionX)

				staticsprite := assets.NewStaticSprite()
				staticsprite.Load("assets/images/bullets/arrow/arrow.png")

				staticsprite.Draw(screen, assets.DrawInput{
					Width:  projectile.Width,
					Height: projectile.Height,
					X:      screenX,
					Y:      screenY,
					Angle:  &angle,
				})
			}
		}
	}

	return
}

func (b *Basic) ActiveProjectiles() []*entities.Projectile {
	return b.Projectiles
}

func (b *Basic) GetPower() float64 {
	return b.Power
}

func (*Basic) DamageType() string {
	return "projectil"
}

func (*Basic) AttackSpeed() float64 {
	return 1.0
}

func (*Basic) GetRadius() float64 {
	return 0.0
}

func (b *Basic) CurrentLevel() int {
	return b.Level
}

func (b *Basic) MaxLevel() bool {
	return b.Level >= 5
}

func (b *Basic) IncreaseLevel() {
	b.Level++
	b.Power += 5
	b.ShootCooldown -= 0.1
}

func (b *Basic) Combat(ci abilityentities.CombatInput) abilityentities.CombatOutput {
	enemy := ci.Enemy
	pp := ci.PlayerPlugin
	enemyGotDamaged := false
	damage := 0.0
	critical := false

	for _, projectil := range b.ActiveProjectiles() {
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

				damage, critical = pp.CalculateDamage(projectil.Power)

				projectil.Active = false
				enemyGotDamaged = true
			}
		}
	}

	return abilityentities.CombatOutput{
		EnemyGotDamaged: enemyGotDamaged,
		Damage:          damage,
		CriticalDamage:  critical,
	}
}
