package weapon

import "github.com/hajimehoshi/ebiten/v2"

type Projectile interface {
	Shoot(x, y float64)
	Print(screen *ebiten.Image)
}
