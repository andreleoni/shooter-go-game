package entities

import (
	"game/internal/assets"
)

type EnemyStats struct {
	MaxHealth float64
	Speed     float64
	Damage    float64
	Size      float64
	Power     float64

	StaticSprite *assets.StaticSprite
}
