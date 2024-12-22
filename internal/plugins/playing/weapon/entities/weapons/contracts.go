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

type Weapon interface {
	AutoShot(deltaTime, x, y float64)
	Shoot(x, y float64)
	Update(WeaponUpdateInput)
	Draw(screen *ebiten.Image, cameraX, cameraY float64)

	ActiveProjectiles() []*entities.Projectile

	GetPower() float64
}
