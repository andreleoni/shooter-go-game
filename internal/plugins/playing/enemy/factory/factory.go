package factory

import (
	"game/internal/plugins/playing/enemy/entities"
	"game/internal/plugins/playing/enemy/templates"
)

func CreateEnemy(enemyType entities.EnemyType, x, y float64) *entities.Enemy {
	template := templates.EnemyTemplates[enemyType]

	return &entities.Enemy{
		X:                                x,
		Y:                                y,
		Width:                            template.Size,
		Height:                           template.Size,
		Active:                           true,
		Type:                             enemyType,
		Stats:                            *template,
		Speed:                            template.Speed,
		Health:                           template.MaxHealth,
		Power:                            template.Power,
		LastAreaDamageDeltaTimeByAbility: map[string]float64{},
	}
}
