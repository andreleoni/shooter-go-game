package protection

import (
	"game/internal/core"
	"game/internal/helpers/collision"
	abilityentities "game/internal/plugins/playing/ability/entities/abilities"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Protection struct {
	plugins *core.PluginManager
	Power   float64
	Radius  float64

	// Shoot cooldown
	ShootTimer    float64
	ShootCooldown float64

	LastDamageDeltaTimeByEnemy map[string]float64

	Level int
}

func New() *Protection {
	return &Protection{
		Power:                      10,
		Radius:                     75,
		Level:                      1,
		LastDamageDeltaTimeByEnemy: make(map[string]float64),
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

func (p *Protection) CurrentLevel() int {
	return p.Level
}

func (p *Protection) MaxLevel() bool {
	return p.Level == 5
}

func (p *Protection) IncreaseLevel() {
	p.Level++
	p.Power += 10
	p.Radius += 10
}

func (p *Protection) Combat(ci abilityentities.CombatInput) abilityentities.CombatOutput {
	enemy := ci.Enemy
	pp := ci.PlayerPlugin
	enemyGotDamaged := false
	damage := 0.0
	critical := false

	playerX, playerY := pp.GetPosition()

	if enemy.Active {
		enemyUUID := enemy.UUID

		if collision.CheckCircle(
			playerX,
			playerY,
			p.GetRadius(),
			enemy.X,
			enemy.Y,
			enemy.Width,
			enemy.Height) {

			lastAreaDamageDeltaTime, exists := p.LastDamageDeltaTimeByEnemy[enemyUUID]
			if !exists {
				lastAreaDamageDeltaTime = 0
			}

			if lastAreaDamageDeltaTime >= p.AttackSpeed() {
				damage, critical = pp.CalculateDamage(p.GetPower())
				enemyGotDamaged = true

				lastAreaDamageDeltaTime = 0
			} else {
				lastAreaDamageDeltaTime += ci.DeltaTime
			}

			p.LastDamageDeltaTimeByEnemy[enemyUUID] = lastAreaDamageDeltaTime
		}
	}

	return abilityentities.CombatOutput{
		EnemyGotDamaged: enemyGotDamaged,
		Damage:          damage,
		CriticalDamage:  critical,
	}
}
