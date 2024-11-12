package collision

import (
	"game/internal/core"
	"game/internal/plugins/bullet"
	"game/internal/plugins/enemy"

	"github.com/hajimehoshi/ebiten/v2"
)

type CollisionPlugin struct {
	kernel       *core.GameKernel
	bulletPlugin *bullet.BulletPlugin
	enemyPlugin  *enemy.EnemyPlugin
}

func NewCollisionPlugin(bulletPlugin *bullet.BulletPlugin, enemyPlugin *enemy.EnemyPlugin) *CollisionPlugin {
	return &CollisionPlugin{
		bulletPlugin: bulletPlugin,
		enemyPlugin:  enemyPlugin,
	}
}

func (cp *CollisionPlugin) ID() string {
	return "CollisionSystem"
}

func (cp *CollisionPlugin) Init(kernel *core.GameKernel) error {
	cp.kernel = kernel
	return nil
}

func (cp *CollisionPlugin) Update() error {
	bullets := cp.bulletPlugin.GetBullets()
	enemies := cp.enemyPlugin.GetEnemies()

	for _, bullet := range bullets {
		if bullet.Active {
			for _, enemy := range enemies {
				if enemy.Active && cp.checkCollision(bullet, enemy) {
					bullet.Active = false
					enemy.Active = false
				}
			}
		}
	}
	return nil
}

func (cp *CollisionPlugin) checkCollision(bullet *bullet.Bullet, enemy *enemy.Enemy) bool {
	return bullet.X < enemy.X+20 && bullet.X+5 > enemy.X && bullet.Y < enemy.Y+20 && bullet.Y+10 > enemy.Y
}

func (cp *CollisionPlugin) Draw(screen *ebiten.Image) {
	// No drawing needed for collision system
}
