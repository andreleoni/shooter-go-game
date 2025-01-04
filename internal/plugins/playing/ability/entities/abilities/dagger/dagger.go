package dagger

import (
	"game/internal/constants"
	"game/internal/core"
	"game/internal/helpers/collision"
	abilityentities "game/internal/plugins/playing/ability/entities/abilities"
	entityabilities "game/internal/plugins/playing/ability/entities/abilities"
	"game/internal/plugins/playing/camera"

	"image/color"
	"math"
	"math/rand/v2"

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
}

type Dagger struct {
	plugins            *core.PluginManager
	Projectiles        []*Projectile
	Power              float64
	ProjectilesByShoot int

	// Shoot cooldown
	ShootTimer    float64
	ShootCooldown float64

	Level int
}

func New() *Dagger {
	return &Dagger{
		Power:              10,
		ShootCooldown:      2.3,
		Level:              1,
		ProjectilesByShoot: 5,
	}
}

func (d *Dagger) SetPluginManager(plugins *core.PluginManager) {
	d.plugins = plugins
}

func (d *Dagger) ID() string {
	return "Dagger"
}

func (d *Dagger) AutoShot(deltaTime, x, y float64) {
	d.ShootTimer += deltaTime

	if d.ShootTimer >= d.ShootCooldown {
		d.Shoot(x, y)
		d.ShootTimer = 0
	}
}

func (d *Dagger) Shoot(x, y float64) {
	for i := 0; i < d.ProjectilesByShoot; i++ {
		angle := rand.Float64() * 2 * math.Pi
		directionX := math.Cos(angle)
		directionY := math.Sin(angle)

		projectile := &Projectile{
			X:          x,
			Y:          y,
			Speed:      300,
			DirectionX: directionX,
			DirectionY: directionY,
			Active:     true,
			Power:      d.Power,
			Height:     5,
			Width:      5,
		}

		d.Projectiles = append(d.Projectiles, projectile)
	}
}

func (d *Dagger) Update(wui entityabilities.AbilityUpdateInput) {
	deltatime := wui.DeltaTime
	d.AutoShot(deltatime, wui.PlayerX, wui.PlayerY)

	for _, p := range d.Projectiles {
		if !p.Active {
			continue
		}

		p.X += p.DirectionX * p.Speed * deltatime
		p.Y += p.DirectionY * p.Speed * deltatime

		cameraPlugin := d.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
		cameraX, cameraY := cameraPlugin.GetPosition()

		screenX := p.X - cameraX
		screenY := p.Y - cameraY
		margin := float64(200)

		if screenX < -margin ||
			screenX > constants.ScreenWidth+margin ||
			screenY < -margin ||
			screenY > constants.ScreenHeight+margin {

			p.Active = false
		}
	}
}

func (d *Dagger) Draw(screen *ebiten.Image, wdi entityabilities.AbilityDrawInput) {
	for _, p := range d.Projectiles {
		if p.Active {
			screenX := p.X - wdi.CameraX
			screenY := p.Y - wdi.CameraY

			vector.DrawFilledRect(
				screen,
				float32(screenX),
				float32(screenY),
				float32(p.Width),
				float32(p.Height),
				color.RGBA{255, 255, 0, 255},
				true)
		}
	}

	return
}

func (d *Dagger) GetPower() float64 {
	return d.Power
}

func (d *Dagger) DamageType() string {
	return "projectil"
}

func (*Dagger) AttackSpeed() float64 {
	return 2.5
}

func (*Dagger) GetRadius() float64 {
	return 0.0
}

func (d *Dagger) CurrentLevel() int {
	return d.Level
}

func (d *Dagger) MaxLevel() bool {
	return d.Level == 5
}

func (d *Dagger) IncreaseLevel() {
	d.Level++
	d.Power += 10
	d.ProjectilesByShoot++
}

func (d *Dagger) Combat(ci abilityentities.CombatInput) abilityentities.CombatOutput {
	enemy := ci.Enemy
	pp := ci.PlayerPlugin
	enemyGotDamaged := false
	damage := 0.0
	critical := false

	for _, projectil := range d.Projectiles {
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
