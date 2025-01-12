package basic

import (
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/helpers/collision"
	"game/internal/plugins"
	abilityentities "game/internal/plugins/playing/ability/entities/abilities"
	"game/internal/plugins/playing/camera"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Projectile struct {
	Active bool
	Power  float64

	X, Y       float64
	Speed      float64
	DirectionX float64
	DirectionY float64

	TargetX float64
	TargetY float64

	Width  float64
	Height float64

	Animation *assets.Animation
}

type Basic struct {
	plugins     *core.PluginManager
	Projectiles []*Projectile
	Power       float64

	// Shoot cooldown
	ShootTimer    float64
	ShootCooldown float64

	Level int

	BaseAnimation *assets.Animation
}

func New() *Basic {

	BaseAnimation := assets.NewAnimation(0.1)
	err := BaseAnimation.LoadFromJSON(
		"assets/images/bullets/fireball/asset.json",
		"assets/images/bullets/fireball/asset.png")

	if err != nil {
		log.Fatal("Failed to load player asset right:", err)
	}

	return &Basic{
		Power:         10,
		ShootCooldown: 1.0,
		Level:         1,
		BaseAnimation: BaseAnimation,
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
		projectile := &Projectile{
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
			Animation:  b.BaseAnimation,
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
			}
		}
	}

	return
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

	for _, projectil := range b.Projectiles {
		if enemy.Active && projectil.Active {
			checkSpriteCollisionInput := collision.CheckSpriteCollisionInput{
				X1:      projectil.X,
				Y1:      projectil.Y,
				Width1:  projectil.Width,
				Height1: projectil.Height,
				X2:      enemy.X,
				Y2:      enemy.Y,
				Width2:  enemy.Width,
				Height2: enemy.Height,
			}

			if collision.CheckSpriteCollision(checkSpriteCollisionInput) {
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
