package abilities

import (
	"game/internal/core"
	"game/internal/plugins/playing/ability/entities"

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

type Ability interface {
	ID() string
	SetPluginManager(plugins *core.PluginManager)
	Shoot(x, y float64)
	Update(wui AbilityUpdateInput)
	Draw(screen *ebiten.Image, wdi AbilityDrawInput)

	ActiveProjectiles() []*entities.Projectile

	GetPower() float64

	DamageType() string

	AttackSpeed() float64
}
