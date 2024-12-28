package entities

import "game/internal/assets"

type Enemy struct {
	X, Y                             float64
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

	RunningAnimationSprite *assets.Animation
	DeathAnimation         *assets.Animation
}

func (e *Enemy) GetBounds() (float64, float64, float64, float64) {
	return e.X, e.Y, e.Width, e.Height
}
