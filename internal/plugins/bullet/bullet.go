// plugins/bullet.go
package bullet

import (
	"game/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Bullet struct {
	X, Y   float64
	Speed  float64
	Active bool
}

type BulletPlugin struct {
	kernel  *core.GameKernel
	bullets []*Bullet
}

func NewBulletPlugin() *BulletPlugin {
	return &BulletPlugin{
		bullets: []*Bullet{},
	}
}

func (bp *BulletPlugin) ID() string {
	return "BulletSystem"
}

func (bp *BulletPlugin) Init(kernel *core.GameKernel) error {
	bp.kernel = kernel
	return nil
}

func (bp *BulletPlugin) Update() error {
	for _, bullet := range bp.bullets {
		if bullet.Active {
			bullet.Y -= bullet.Speed * bp.kernel.DeltaTime
			if bullet.Y < 0 {
				bullet.Active = false
			}
		}
	}
	return nil
}

func (bp *BulletPlugin) Draw(screen *ebiten.Image) {
	for _, bullet := range bp.bullets {
		if bullet.Active {
			ebitenutil.DrawRect(screen, bullet.X, bullet.Y, 5, 10, color.RGBA{255, 0, 0, 255})
		}
	}
}

func (bp *BulletPlugin) Shoot(x, y float64) {
	bp.bullets = append(bp.bullets, &Bullet{X: x, Y: y, Speed: 300, Active: true})
}

func (bp *BulletPlugin) GetBullets() []*Bullet {
	return bp.bullets
}
