package combat

import (
	"game/internal/core"
	"game/internal/plugins/bullet"
	"game/internal/plugins/enemy"

	"github.com/hajimehoshi/ebiten/v2"
)

type CombatPlugin struct {
	kernel       *core.GameKernel
	bulletPlugin *bullet.BulletPlugin
	enemyPlugin  *enemy.EnemyPlugin
}

func NewCombatPlugin(bulletPlugin *bullet.BulletPlugin, enemyPlugin *enemy.EnemyPlugin) *CombatPlugin {
	return &CombatPlugin{
		bulletPlugin: bulletPlugin,
		enemyPlugin:  enemyPlugin,
	}
}

func (cp *CombatPlugin) ID() string {
	return "CombatSystem"
}

func (cp *CombatPlugin) Init(kernel *core.GameKernel) error {
	cp.kernel = kernel
	return nil
}

func (cp *CombatPlugin) Draw(*ebiten.Image) {
	return
}

func (cp *CombatPlugin) Update() error {
	bullets := cp.bulletPlugin.GetBullets()
	enemies := cp.enemyPlugin.GetEnemies()

	for _, bullet := range bullets {
		if bullet.Active {
			for _, enemy := range enemies {
				if enemy.Active {
					// Check if bullet hits enemy
					if checkCollision(bullet.X, bullet.Y, 5, 10, enemy.X, enemy.Y, 20, 20) {
						bullet.Active = false
						enemy.Active = false
						// Here we could add effects, sounds, score etc
					}
				}
			}
		}
	}
	return nil
}

func checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}
