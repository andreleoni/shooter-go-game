package weapons

import (
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/playing/weapon/entities"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Dagger struct {
	plugins     *core.PluginManager
	Projectiles []*entities.Projectile
	Power       float64

	// Shoot cooldown
	ShootTimer    float64
	ShootCooldown float64
}

func NewDagger() *Dagger {
	return &Dagger{
		Power:         10,
		ShootCooldown: 2.3,
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
	for i := 0; i < 10; i++ {
		angle := rand.Float64() * 2 * math.Pi
		directionX := math.Cos(angle)
		directionY := math.Sin(angle)

		projectile := &entities.Projectile{
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

func (d *Dagger) Update(wui WeaponUpdateInput) {
	deltatime := wui.DeltaTime
	d.AutoShot(deltatime, wui.PlayerX, wui.PlayerY)

	for _, weapon := range d.Projectiles {
		if !weapon.Active {
			continue
		}

		weapon.X += weapon.DirectionX * weapon.Speed * deltatime
		weapon.Y += weapon.DirectionY * weapon.Speed * deltatime

		// Deactivate if off screen
		if weapon.X < 0 || weapon.X > constants.WorldWidth ||
			weapon.Y < 0 || weapon.Y > constants.WorldHeight {
			weapon.Active = false
		}
	}
}

func (d *Dagger) Draw(screen *ebiten.Image, wdi WeaponDrawInput) {
	for _, weapon := range d.Projectiles {
		if weapon.Active {
			screenX := weapon.X - wdi.CameraX
			screenY := weapon.Y - wdi.CameraY

			vector.DrawFilledRect(
				screen,
				float32(screenX),
				float32(screenY),
				float32(weapon.Width),
				float32(weapon.Height),
				color.RGBA{255, 255, 0, 255},
				true)
		}
	}

	return
}

func (d *Dagger) ActiveProjectiles() []*entities.Projectile {
	return d.Projectiles
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
