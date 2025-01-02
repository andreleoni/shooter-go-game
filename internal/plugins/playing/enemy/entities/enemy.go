package entities

import (
	"game/internal/assets"
	"math"
)

type Enemy struct {
	UUID string

	Name string

	X, Y float64

	Width                            float64
	Height                           float64
	Active                           bool
	Speed                            float64
	Type                             EnemyType
	Template                         EnemyTemplate
	Health                           float64
	MaxHealth                        float64
	Power                            float64
	LastPlayerDamageTime             float64
	LastAreaDamageDeltaTimeByAbility map[string]float64
	DamageFlashTime                  float64

	RunningRightAnimationSprite *assets.Animation
	RunningLeftAnimationSprite  *assets.Animation

	RunningAnimationTime float64

	DeathAnimation *assets.Animation

	VelocityX float64
	VelocityY float64
	StuckTime float64

	AttackCooldown  float64
	AttackRange     float64
	ProjectileSpeed float64
	Projectiles     []*Projectile
}

type Projectile struct {
	X, Y          float64
	Width, Height float64
	Speed         float64
	DirectionX    float64
	DirectionY    float64
	Active        bool
	Damage        float64
}

func (e *Enemy) GetBounds() (float64, float64, float64, float64) {
	return e.X, e.Y, e.Width, e.Height
}

func (e *Enemy) IsEnemyMovingRight(playerX float64) bool {
	return e.X < playerX
}

func (re *Enemy) Shoot(playerX, playerY float64) {
	dx := playerX - re.X
	dy := playerY - re.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Normalize direction
	if distance > 0 {
		dx /= distance
		dy /= distance
	}

	projectile := &Projectile{
		X:          re.X,
		Y:          re.Y,
		Width:      8,
		Height:     8,
		Speed:      re.ProjectileSpeed,
		DirectionX: dx,
		DirectionY: dy,
		Active:     true,
		Damage:     10,
	}

	re.Projectiles = append(re.Projectiles, projectile)
}
