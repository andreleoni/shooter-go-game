package protection

import (
	"game/internal/core"
	"game/internal/plugins/playing/ability/entities"
	abilityentities "game/internal/plugins/playing/ability/entities/abilities"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Protection struct {
	plugins     *core.PluginManager
	Projectiles []*entities.Projectile
	Power       float64
	Radius      float64

	// Shoot cooldown
	ShootTimer    float64
	ShootCooldown float64
}

func New() *Protection {
	return &Protection{
		Power:  10,
		Radius: 75,
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

func (p *Protection) Update(wui abilityentities.AbilityUpdateInput) {
}

func (p *Protection) Draw(screen *ebiten.Image, wdi abilityentities.AbilityDrawInput) {
	screenX := wdi.PlayerX - wdi.CameraX
	screenY := wdi.PlayerY - wdi.CameraY

	circleX := screenX
	circleY := screenY

	vector.DrawFilledCircle(
		screen,
		float32(circleX),
		float32(circleY),
		float32(p.Radius),
		color.RGBA{111, 222, 111, 2},
		true)

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

func (p *Protection) GetRadius() float64 {
	return p.Radius
}
