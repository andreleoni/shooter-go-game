package protection

import (
	"game/internal/constants"
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
		Radius:                     200,
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

	vector.DrawFilledCircle(
		screen,
		float32(screenX),
		float32(screenY),
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

	if enemy.Active {
		enemyUUID := enemy.UUID

		screenCenterX := float64(constants.ScreenWidth) / 2
		screenCenterY := float64(constants.ScreenHeight) / 2

		circleCenterX := screenCenterX + ci.CameraX - p.GetRadius()
		circleCenterY := screenCenterY + ci.CameraY - p.GetRadius()

		enemyCenterX := enemy.X + enemy.Width/2
		enemyCenterY := enemy.Y + enemy.Height/2

		checkSpriteCollisionInput := collision.CheckSpriteCollisionInput{
			X1:      circleCenterX,
			Y1:      circleCenterY,
			Width1:  p.GetRadius() * 2,
			Height1: p.GetRadius() * 2,
			X2:      enemyCenterX,
			Y2:      enemyCenterY,
			Width2:  enemy.Width,
			Height2: enemy.Height,
		}

		if collision.CheckSpriteCollision(checkSpriteCollisionInput) {
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
