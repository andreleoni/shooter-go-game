package enemy

import (
	"fmt"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/camera"
	entity "game/internal/plugins/enemy/entities"
	"game/internal/plugins/enemy/factory"
	"game/internal/plugins/enemy/templates"
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
	enemies      []*entity.Enemy
	spawnTimer   float64
	playerPlugin *player.PlayerPlugin
}

func NewEnemyPlugin(playerPlugin *player.PlayerPlugin) *EnemyPlugin {
	return &EnemyPlugin{
		enemies:      []*entity.Enemy{},
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

	for i, enemy := range ep.enemies {
		if enemy.Active {
			// Store current position
			oldX, oldY := enemy.X, enemy.Y

			// Calculate desired velocity towards player
			dx := playerX - enemy.X
			dy := playerY - enemy.Y
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance > 0 {
				dx /= distance
				dy /= distance

				desiredX := dx * enemy.Speed * ep.kernel.DeltaTime
				desiredY := dy * enemy.Speed * ep.kernel.DeltaTime

				// Check for obstacles and steer around them
				if !obstaclePlugin.CheckCollisionRect(enemy.X+desiredX, enemy.Y+desiredY, enemy.Width, enemy.Height) {
					enemy.X += desiredX
					enemy.Y += desiredY
				} else {
					// Enhanced obstacle avoidance by trying multiple directions
					directions := [][2]float64{
						{-desiredY, desiredX},  // Left
						{desiredY, -desiredX},  // Right
						{-desiredX, -desiredY}, // Backward
					}

					moved := false

					for _, dir := range directions {
						if !obstaclePlugin.CheckCollisionRect(enemy.X+dir[0], enemy.Y+dir[1], enemy.Width, enemy.Height) {
							enemy.X += dir[0]
							enemy.Y += dir[1]

							moved = true

							break
						}
					}

					// If no direction is found, revert to old position
					if !moved {
						enemy.X = oldX
						enemy.Y = oldY
					}
				}

				// Check for collisions with other enemies
				for j, otherEnemy := range ep.enemies {
					if i != j && otherEnemy.Active {
						if math.Abs(enemy.X-otherEnemy.X) < enemy.Width && math.Abs(enemy.Y-otherEnemy.Y) < enemy.Height {
							// Adjust position to avoid collision
							enemy.X = oldX
							enemy.Y = oldY
							break
						}
					}
				}
			}
		}
	}
	return nil
}

func (ep *EnemyPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := ep.kernel.PluginManager.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for _, enemy := range ep.enemies {
		if enemy.Active {
			// Draw enemy relative to camera position
			screenX := enemy.X - cameraX
			screenY := enemy.Y - cameraY

			// Only draw if on screen
			if screenX >= -enemy.Width && screenX <= constants.ScreenWidth+enemy.Width &&
				screenY >= -enemy.Height && screenY <= constants.ScreenHeight+enemy.Height {

				ebitenutil.DrawRect(
					screen,
					screenX,
					screenY,
					enemy.Width,
					enemy.Height,
					color.RGBA{255, 0, 0, 255},
				)
			}
		}
	}
}

func (ep *EnemyPlugin) Spawn(x, y float64) {
	fmt.Println("Spawn enemy at", x, y)
	ep.enemies = append(ep.enemies, factory.CreateEnemy(templates.TankEnemy, x, y))
}

func (ep *EnemyPlugin) GetEnemies() []*entity.Enemy {
	return ep.enemies
}
