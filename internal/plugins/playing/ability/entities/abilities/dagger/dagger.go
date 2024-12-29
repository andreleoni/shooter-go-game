package dagger

import (
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/playing/ability/entities"
	entityabilities "game/internal/plugins/playing/ability/entities/abilities"
	"game/internal/plugins/playing/camera"

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

func New() *Dagger {
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
		// Check if projectile is too far from camera view
		screenX := p.X - cameraX
		screenY := p.Y - cameraY
		margin := float64(200) // Larger margin before deactivating

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
