package weapons

import (
	"game/internal/plugins/playing/weapon/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

type WeaponUpdateInput struct {
	DeltaTime float64

	PlayerX float64
	PlayerY float64

	CameraX float64
	CameraY float64
}

type WeaponDrawInput struct {
	PlayerX float64
	PlayerY float64

	CameraX float64
	CameraY float64
}

type Weapon interface {
	ID() string
	AutoShot(deltaTime, x, y float64)
	Shoot(x, y float64)
	Update(wui WeaponUpdateInput)
	Draw(screen *ebiten.Image, wdi WeaponDrawInput)

	ActiveProjectiles() []*entities.Projectile

	GetPower() float64

	DamageType() string

	AttackSpeed() float64
}
