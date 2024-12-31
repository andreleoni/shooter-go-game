package entities

import "game/internal/assets"

type Character struct {
	Name           string
	ID             string
	Speed          float64
	Health         float64
	Ability        string
	Armor          float64
	DamagePercent  float64
	CriticalChance float64

	HealthRegenRate  float64
	HealthRegenDelay float64
	HealthRegenTimer float64

	animation *assets.Animation
}
