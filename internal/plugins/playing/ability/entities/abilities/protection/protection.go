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

	// Shoot cooldown
	ShootTimer    float64
	ShootCooldown float64
}

func New() *Protection {
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
		200,
		color.RGBA{111, 222, 111, 255},
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
