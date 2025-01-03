package factory

import (
	"game/internal/assets"
	"game/internal/plugins/playing/enemy/entities"
	"game/internal/plugins/playing/enemy/templates"

	"github.com/google/uuid"
)

func CreateEnemy(enemyType entities.EnemyType, x, y float64) *entities.Enemy {
	template := templates.EnemyTemplates[enemyType]

	randomUUID := uuid.NewString()

	return &entities.Enemy{
		Name:                             template.Name,
		UUID:                             randomUUID,
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

		AttackCooldown:  2.0, // Seconds between attacks
		AttackRange:     800, // Range to start shooting
		ProjectileSpeed: 200,
	}
}

func dupAnimation(a *assets.Animation) *assets.Animation {
	val := *a

	return &val
}
