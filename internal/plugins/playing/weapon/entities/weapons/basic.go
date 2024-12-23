package weapons

import (
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins"
	"game/internal/plugins/playing/weapon/entities"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Basic struct {
	plugins     *core.PluginManager
	Projectiles []*entities.Projectile
	Power       float64

	// Shoot cooldown
	ShootTimer    float64
	ShootCooldown float64
}

func NewBasic(plugins *core.PluginManager) *Basic {
	return &Basic{
		Power:         10,
		ShootCooldown: 1.0,
		plugins:       plugins,
	}
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

func (b *Basic) Update(wui WeaponUpdateInput) {
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

			// Deactivate if off screen
			if projectile.X < 0 ||
				projectile.X > constants.WorldHeight ||
				projectile.Y < 0 ||
				projectile.Y > constants.WorldWidth {

				projectile.Active = false
			}
		}
	}
}

func (b *Basic) Draw(screen *ebiten.Image, wdi WeaponDrawInput) {
	for _, projectile := range b.Projectiles {
		if projectile.Active {
			// Draw bullet relative to camera position
			screenX := projectile.X - wdi.CameraX
			screenY := projectile.Y - wdi.CameraY

			// Only draw if on screen
			if screenX >= -5 && screenX <= constants.ScreenWidth+5 &&
				screenY >= -5 && screenY <= constants.ScreenHeight+5 {

				ebitenutil.DrawRect(screen, screenX, screenY, projectile.Width, projectile.Height, color.RGBA{200, 255, 0, 255})

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
