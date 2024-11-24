package entities

import (
	"game/internal/assets"
	"image/color"
)

type EnemyStats struct {
	MaxHealth float64
	Speed     float64
	Damage    float64
	Size      float64
	Color     color.RGBA
	Power     float64

	StaticSprite *assets.StaticSprite
}
