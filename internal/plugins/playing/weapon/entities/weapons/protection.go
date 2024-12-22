package weapons

import (
	"game/internal/core"
	"game/internal/plugins"
	"game/internal/plugins/playing/weapon/entities"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Protection struct {
	plugins     *core.PluginManager
	Projectiles []*entities.Projectile
	Power       float64

	// Shoot cooldown
	ShootTimer    float64
	ShootCooldown float64
}

func NewProtection() *Protection {
	return &Protection{
		Power: 10,
	}
}

func (p *Protection) ID() string {
	return "Protection"
}

func (p *Protection) Shoot(deltaTime, x, y float64) {
	return
}

func (p *Protection) Update(deltaTime, px, py float64) {
}

func (p *Protection) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	playerPlugin := p.plugins.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)
	playerX, playerY := playerPlugin.GetPosition()

	screenX := playerX - cameraX
	screenY := playerY - cameraY

	circleX := screenX
	circleY := screenY

	ebitenutil.DrawCircle(screen, circleX, circleY, 50, color.RGBA{111, 222, 111, 255})

	return
}

func (d *Protection) ActiveProjectiles() []*entities.Projectile {
	return d.Projectiles
}

func (d *Protection) GetPower() float64 {
	return d.Power
}
