package templates

import (
	"game/internal/plugins/playing/enemy/entities"
)

const (
	BasicEnemy entities.EnemyType = iota
	FastEnemy
	TankEnemy
	RangedEnemy
)

var EnemyTemplates = map[entities.EnemyType]*entities.EnemyStats{
	BasicEnemy: {
		MaxHealth: 20,
		Speed:     100,
		Damage:    10,
		Size:      40,
		Power:     10,
	},
	FastEnemy: {
		MaxHealth: 10,
		Speed:     200,
		Damage:    5,
		Size:      100,
		Power:     5,
	},
	TankEnemy: {
		MaxHealth: 50,
		Speed:     50,
		Damage:    20,
		Size:      60,
		Power:     20,
	},
	RangedEnemy: {
		MaxHealth: 30,
		Speed:     75,
		Damage:    15,
		Size:      18,
		Power:     15,
	},
}
