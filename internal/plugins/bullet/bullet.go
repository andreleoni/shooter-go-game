// plugins/bullet.go
package bullet

import (
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins"
	entities "game/internal/plugins/bullet/entities"
	"game/internal/plugins/camera"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type BulletPlugin struct {
	kernel  *core.GameKernel
	bullets []*entities.Bullet

	sprite *assets.StaticSprite
}

func NewBulletPlugin() *BulletPlugin {
	return &BulletPlugin{
		bullets: []*entities.Bullet{},
	}
}

func (bp *BulletPlugin) ID() string {
	return "BulletSystem"
}

func (bp *BulletPlugin) Init(kernel *core.GameKernel) error {
	bp.kernel = kernel

	bp.sprite = assets.NewStaticSprite()
	bp.sprite.Load("assets/images/bullets/arrow/arrow.png")

	return nil
}

func (bp *BulletPlugin) Update() error {
	for _, bullet := range bp.bullets {
		if bullet.Active {
			bullet.MoveTowardsTarget(bp.kernel.DeltaTime)

			// Deactivate if off screen
			if bullet.X < 0 ||
				bullet.X > constants.WorldHeight ||
				bullet.Y < 0 ||
				bullet.Y > constants.WorldWidth {

				bullet.Active = false
			}
		}
	}
	return nil
}

func (bp *BulletPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := bp.kernel.PluginManager.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for _, bullet := range bp.bullets {
		if bullet.Active {
			// Draw bullet relative to camera position
			screenX := bullet.X - cameraX
			screenY := bullet.Y - cameraY

			// Only draw if on screen
			if screenX >= -5 && screenX <= constants.ScreenWidth+5 &&
				screenY >= -5 && screenY <= constants.ScreenHeight+5 {

				angle := math.Atan2(bullet.DirectionY, bullet.DirectionX)

				bp.sprite.DrawAngle(screen, screenX, screenY, angle)
			}
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

		// Calcular direção
		dx := closestEnemy.X - x
		dy := closestEnemy.Y - y

		distance := math.Sqrt(dx*dx + dy*dy)

		dirX := dx / distance
		dirY := dy / distance

		// Create bullet targeting closest enemy
		bullet := &entities.Bullet{
			X:          x,
			Y:          y,
			Speed:      300,
			Active:     true,
			Power:      50,
			TargetX:    closestEnemy.X,
			TargetY:    closestEnemy.Y,
			DirectionX: dirX,
			DirectionY: dirY,
		}

		bp.bullets = append(bp.bullets, bullet)
	}
}

func (bp *BulletPlugin) GetBullets() []*entities.Bullet {
	return bp.bullets
}
