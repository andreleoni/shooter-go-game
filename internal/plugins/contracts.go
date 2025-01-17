package plugins

import (
	"game/internal/core"
	"game/internal/plugins/playing/enemy/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyPlugin interface {
	ID() string
	Init(kernel *core.GameKernel) error
	Update() error
	Draw(screen *ebiten.Image)
	Spawn()
	GetEnemies() []*entities.Enemy
	ApplyDamage(enemy *entities.Enemy, damage float64, isCriticalDamage bool)
	GetGlobalProjectiles() []*entities.Projectile
}

type PlayerPlugin interface {
	ID() string
	Init(kernel *core.GameKernel) error
	Update() error
	Draw(screen *ebiten.Image)
	GetPosition() (float64, float64)
	DecreaseHealth(float64)

	GetSize() (float64, float64)
	GetLevel() float64
	GetExperience() float64

	GetHealth() float64
	GetMaxHealth() float64
	GetArmor() float64
	GetDamagePercent() float64
	GetCriticalChance() float64
	GetSpeed() float64
	GetHealthRegenRate() float64
	GetHealthRegenDelay() float64
	GetAdditionalDamagePercent() float64
	GetCollectionRadius() float64

	NextLevelPercentage() float64

	ApplyDamage(damage float64)
	CalculateDamage(baseDamage float64) (float64, bool)

	GetDashTimer() float64

	GetNextLevelExperience() int

	AddExperience(amount int)
}
