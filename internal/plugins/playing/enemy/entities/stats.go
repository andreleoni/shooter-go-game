package entities

import (
	"game/internal/assets"
)

type EnemyTemplate struct {
	Name      string
	MaxHealth float64
	Speed     float64
	Damage    float64
	Size      float64
	Power     float64

	RunningAnimationSprite *assets.Animation
	DeathAnimation         *assets.Animation
}
