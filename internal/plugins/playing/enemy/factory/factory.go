package factory

import (
	"game/internal/assets"
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
		Template:                         *template,
		Speed:                            template.Speed,
		Health:                           template.MaxHealth,
		Power:                            template.Power,
		LastAreaDamageDeltaTimeByAbility: map[string]float64{},
		MaxHealth:                        template.MaxHealth,
		RunningRightAnimationSprite:      dupAnimation(template.RunningRightAnimationSprite),
		RunningLeftAnimationSprite:       dupAnimation(template.RunningLeftAnimationSprite),
		DeathAnimation:                   dupAnimation(template.DeathAnimation),
	}
}

func dupAnimation(a *assets.Animation) *assets.Animation {
	val := *a

	return &val
}
