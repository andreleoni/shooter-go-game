package templates

import (
	"game/internal/plugins/enemy/entities"
	"image/color"
)

const (
	BasicEnemy  entities.EnemyType = "basic"
	FastEnemy   entities.EnemyType = "fast"
	TankEnemy   entities.EnemyType = "tank"
	RangedEnemy entities.EnemyType = "ranged"
)

var EnemyTemplates = map[entities.EnemyType]entities.EnemyStats{
	BasicEnemy: {
		MaxHealth: 100,
		Speed:     100,
		Damage:    10,
		Size:      20,
		Color:     color.RGBA{255, 0, 0, 255},
		Power:     10,
	},
	FastEnemy: {
		MaxHealth: 50,
		Speed:     200,
		Damage:    5,
		Size:      15,
		Color:     color.RGBA{0, 255, 0, 255},
		Power:     5,
	},
	TankEnemy: {
		MaxHealth: 200,
		Speed:     50,
		Damage:    20,
		Size:      30,
		Color:     color.RGBA{0, 0, 255, 255},
		Power:     20,
	},
	RangedEnemy: {
		MaxHealth: 75,
		Speed:     75,
		Damage:    15,
		Size:      18,
		Color:     color.RGBA{255, 255, 0, 255},
		Power:     15,
	},
}
