// plugins/bullet.go
package bullet

import (
	"game/internal/core"
	"game/internal/plugins"
	bulletentities "game/internal/plugins/bullet/entities"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type BulletPlugin struct {
	kernel  *core.GameKernel
	bullets []*bulletentities.Bullet
}

func NewBulletPlugin() *BulletPlugin {
	return &BulletPlugin{
		bullets: []*bulletentities.Bullet{},
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
			bullet.MoveTowardsTarget(bp.kernel.DeltaTime)

			// Deactivate if off screen
			if bullet.X < 0 ||
				bullet.X > 800 ||
				bullet.Y < 0 ||
				bullet.Y > 600 {

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
	// Get enemy plugin to find closest enemy
	enemyPlugin := bp.kernel.PluginManager.GetPlugin("EnemySystem").(plugins.EnemyPlugin)

	enemies := enemyPlugin.GetEnemies()

	if len(enemies) > 0 {
		// Find closest enemy
		closestEnemy := enemies[0]
		closestDist := math.MaxFloat64

		for _, enemy := range enemies {
			if !enemy.Active {
				continue
			}

			dx := enemy.X - x
			dy := enemy.Y - y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist < closestDist {
				closestDist = dist
				closestEnemy = enemy
			}
		}

		// Create bullet targeting closest enemy
		bullet := &bulletentities.Bullet{
			X:       x,
			Y:       y,
			Speed:   300,
			Active:  true,
			TargetX: closestEnemy.X,
			TargetY: closestEnemy.Y,
		}

		bp.bullets = append(bp.bullets, bullet)
	}
}

func (bp *BulletPlugin) GetBullets() []*bulletentities.Bullet {
	return bp.bullets
}
