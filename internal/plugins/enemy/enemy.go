package enemy

import (
	"game/internal/core"
	"game/internal/plugins/obstacle"
	"game/internal/plugins/player"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EnemyPlugin struct {
	kernel       *core.GameKernel
	enemies      []*Enemy
	spawnTimer   float64
	playerPlugin *player.PlayerPlugin
}

func NewEnemyPlugin(playerPlugin *player.PlayerPlugin) *EnemyPlugin {
	return &EnemyPlugin{
		enemies:      []*Enemy{},
		spawnTimer:   0,
		playerPlugin: playerPlugin,
	}
}

func (ep *EnemyPlugin) ID() string {
	return "EnemySystem"
}

func (ep *EnemyPlugin) Init(kernel *core.GameKernel) error {
	ep.kernel = kernel
	return nil
}

func (ep *EnemyPlugin) Update() error {
	ep.spawnTimer += ep.kernel.DeltaTime
	if ep.spawnTimer >= 1.0 {
		ep.Spawn(rand.Float64()*800, 0)
		ep.spawnTimer = 0
	}

	playerX, playerY := ep.playerPlugin.GetPosition()
	obstaclePlugin := ep.kernel.PluginManager.GetPlugin("ObstacleSystem").(*obstacle.ObstaclePlugin)

	for _, enemy := range ep.enemies {
		if enemy.Active {
			// Store current position
			oldX, oldY := enemy.X, enemy.Y

			// Calculate movement
			dx := playerX - enemy.X
			dy := playerY - enemy.Y
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance > 0 {
				dx /= distance
				dy /= distance

				newX := enemy.X + dx*enemy.Speed*ep.kernel.DeltaTime
				newY := enemy.Y + dy*enemy.Speed*ep.kernel.DeltaTime

				// Check collision with enemy size (20x20)
				if !obstaclePlugin.CheckCollisionRect(newX, newY, 20, 20) {
					enemy.X = newX
					enemy.Y = newY
				} else {
					// Revert position if collision detected
					enemy.X = oldX
					enemy.Y = oldY
				}
			}
		}
	}
	return nil
}
func (ep *EnemyPlugin) Draw(screen *ebiten.Image) {
	for _, enemy := range ep.enemies {
		if enemy.Active {
			ebitenutil.DrawRect(screen, enemy.X, enemy.Y, 20, 20, color.RGBA{255, 0, 0, 255})
		}
	}
}

func (ep *EnemyPlugin) Spawn(x, y float64) {
	ep.enemies = append(ep.enemies, &Enemy{X: x, Y: y, Speed: 100, Active: true})
}

func (ep *EnemyPlugin) GetEnemies() []*Enemy {
	return ep.enemies
}
