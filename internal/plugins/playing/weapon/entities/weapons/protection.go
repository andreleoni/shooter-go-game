package weapons

import (
	"game/internal/core"
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

func (p *Protection) SetPluginManager(plugins *core.PluginManager) {
	p.plugins = plugins
}

func (p *Protection) ID() string {
	return "Protection"
}

func (p *Protection) Shoot(x, y float64) {
	return
}

func (p *Protection) Update(wui WeaponUpdateInput) {
}

func (p *Protection) Draw(screen *ebiten.Image, wdi WeaponDrawInput) {
	screenX := wdi.PlayerX - wdi.CameraX
	screenY := wdi.PlayerY - wdi.CameraY

	circleX := screenX
	circleY := screenY

	ebitenutil.DrawCircle(screen, circleX, circleY, 200, color.RGBA{111, 222, 111, 255})

	return
}

func (d *Protection) ActiveProjectiles() []*entities.Projectile {
	return d.Projectiles
}

func (d *Protection) GetPower() float64 {
	return d.Power
}

func (*Protection) DamageType() string {
	return "area"
}

func (*Protection) AttackSpeed() float64 {
	return 0.5
}

func (*Protection) AutoShot(deltaTime, x, y float64) {
}
