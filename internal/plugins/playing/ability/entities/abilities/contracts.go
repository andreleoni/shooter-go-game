package abilities

import (
	"game/internal/core"
	"game/internal/plugins"
	enemyentities "game/internal/plugins/playing/enemy/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

type AbilityUpdateInput struct {
	DeltaTime float64

	PlayerX float64
	PlayerY float64

	CameraX float64
	CameraY float64
}

type AbilityDrawInput struct {
	PlayerX float64
	PlayerY float64

	CameraX float64
	CameraY float64
}

type CombatInput struct {
	DeltaTime float64

	PlayerPlugin plugins.PlayerPlugin
	EnemyPlugin  plugins.EnemyPlugin

	Enemy *enemyentities.Enemy
}

type CombatOutput struct {
	Damage          float64
	CriticalDamage  bool
	EnemyGotDamaged bool
}

type Ability interface {
	ID() string
	SetPluginManager(plugins *core.PluginManager)
	Shoot(x, y float64)
	Update(wui AbilityUpdateInput)
	Draw(screen *ebiten.Image, wdi AbilityDrawInput)

	GetPower() float64

	DamageType() string

	Combat(ci CombatInput) CombatOutput

	AttackSpeed() float64

	GetRadius() float64

	CurrentLevel() int
	MaxLevel() bool
	IncreaseLevel()
}
