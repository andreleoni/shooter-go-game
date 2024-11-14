package entities

import "image/color"

type EnemyStats struct {
	MaxHealth float64
	Speed     float64
	Damage    float64
	Size      float64
	Color     color.RGBA
}
