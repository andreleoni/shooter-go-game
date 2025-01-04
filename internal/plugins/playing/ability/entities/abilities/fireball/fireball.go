package fireball

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

	Radius float64

	EnemiesDamaged map[string]bool

	Animation *assets.Animation
}

type Ability struct {
	plugins     *core.PluginManager
	Projectiles []*Projectile
	Power       float64

	// Shoot cooldown
	ShootTimer    float64
	ShootCooldown float64

	Level int

	FireballAnimation *assets.Animation
}

func New() *Ability {
	fireballAnimation := assets.NewAnimation(0.1)
	err := fireballAnimation.LoadFromJSON(
		"assets/images/bullets/fireball/asset.json",
		"assets/images/bullets/fireball/asset.png")

	if err != nil {
		log.Fatal("Failed to load player asset right:", err)
	}

	return &Ability{
		Power:             10,
		ShootCooldown:     1.0,
		Level:             1,
		FireballAnimation: fireballAnimation,
	}
}

func (b *Ability) SetPluginManager(plugins *core.PluginManager) {
	b.plugins = plugins
}

func (b *Ability) ID() string {
	return "Fireball"
}

func (b *Ability) AutoShot(deltaTime, x, y float64) {
	b.ShootTimer += deltaTime

	if b.ShootTimer >= b.ShootCooldown {
		b.Shoot(x, y)
		b.ShootTimer = 0
	}
}

func (b *Ability) Shoot(x, y float64) {
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
		projectile := &Projectile{
			X:              x,
			Y:              y,
			Speed:          500,
			Active:         true,
			Power:          b.Power,
			DirectionX:     dirX,
			DirectionY:     dirY,
			Radius:         50,
			EnemiesDamaged: make(map[string]bool),
			Animation:      b.FireballAnimation,
		}

		b.Projectiles = append(b.Projectiles, projectile)
	}
}

func (b *Ability) Update(wui abilityentities.AbilityUpdateInput) {
	b.AutoShot(wui.DeltaTime, wui.PlayerX, wui.PlayerY)

	b.FireballAnimation.Update(wui.DeltaTime)

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

func (b *Ability) Draw(screen *ebiten.Image, wdi abilityentities.AbilityDrawInput) {
	for _, projectile := range b.Projectiles {
		if projectile.Active {
			// Draw bullet relative to camera position
			screenX := projectile.X - wdi.CameraX
			screenY := projectile.Y - wdi.CameraY

			// Only draw if on screen
			if screenX >= -5 && screenX <= constants.ScreenWidth+5 &&
				screenY >= -5 && screenY <= constants.ScreenHeight+5 {

				vector.DrawFilledCircle(
					screen,
					float32(screenX),
					float32(screenY),
					float32(projectile.Radius),
					color.RGBA{200, 255, 0, 255},
					true)

				squareSize := projectile.Radius * 2
				drawInput := assets.DrawInput{
					Width:  squareSize,
					Height: squareSize,
					X:      screenX - squareSize/2,
					Y:      screenY - squareSize/2,
				}

				projectile.Animation.Draw(screen, drawInput)
			}
		}
	}

	return
}

func (b *Ability) GetPower() float64 {
	return b.Power
}

func (*Ability) DamageType() string {
	return "projectil"
}

func (*Ability) AttackSpeed() float64 {
	return 1.0
}

func (*Ability) GetRadius() float64 {
	return 0.0
}

func (b *Ability) CurrentLevel() int {
	return b.Level
}

func (b *Ability) MaxLevel() bool {
	return b.Level >= 5
}

func (b *Ability) IncreaseLevel() {
	b.Level++
	b.Power += 5
	b.ShootCooldown -= 0.1
}

func (b *Ability) Combat(ci abilityentities.CombatInput) abilityentities.CombatOutput {
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
				Width1:  projectil.Radius * 2,
				Height1: projectil.Radius * 2,
				X2:      enemy.X,
				Y2:      enemy.Y,
				Width2:  enemy.Width,
				Height2: enemy.Height,
			}

			if collision.CheckSpriteCollision(checkSpriteCollisionInput) {
				damage, critical = pp.CalculateDamage(projectil.Power)

				projectil.EnemiesDamaged[enemy.UUID] = true
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
