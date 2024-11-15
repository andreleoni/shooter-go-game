// internal/plugins/obstacle/obstacle.go
package obstacle

import (
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/camera"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Obstacle struct {
	X, Y          float64
	Width, Height float64
	Active        bool
}

type ObstaclePlugin struct {
	kernel     *core.GameKernel
	obstacles  []*Obstacle
	spawnTimer float64
}

func NewObstaclePlugin() *ObstaclePlugin {
	return &ObstaclePlugin{
		obstacles:  make([]*Obstacle, 0),
		spawnTimer: 0,
	}
}

func (op *ObstaclePlugin) ID() string {
	return "ObstacleSystem"
}

func (op *ObstaclePlugin) Init(kernel *core.GameKernel) error {
	op.kernel = kernel

	op.SpawnRandomObstacle()
	op.SpawnRandomObstacle()
	op.SpawnRandomObstacle()
	op.SpawnRandomObstacle()
	op.SpawnRandomObstacle()
	op.SpawnRandomObstacle()
	op.SpawnRandomObstacle()
	op.SpawnRandomObstacle()
	op.SpawnRandomObstacle()
	op.SpawnRandomObstacle()
	op.SpawnRandomObstacle()

	return nil
}

func (op *ObstaclePlugin) Update() error {
	return nil
}

func (op *ObstaclePlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := op.kernel.PluginManager.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for _, obstacle := range op.obstacles {
		// Draw obstacle relative to camera position
		screenX := obstacle.X - cameraX
		screenY := obstacle.Y - cameraY

		// Only draw if on screen
		if screenX >= -obstacle.Width && screenX <= constants.ScreenWidth+obstacle.Width &&
			screenY >= -obstacle.Height && screenY <= constants.ScreenHeight+obstacle.Height {
			ebitenutil.DrawRect(
				screen,
				screenX,
				screenY,
				obstacle.Width,
				obstacle.Height,
				color.RGBA{0, 255, 0, 255},
			)
		}
	}
}
func (op *ObstaclePlugin) SpawnRandomObstacle() {
	width := 30.0 + rand.Float64()*50.0  // Random width between 30-80
	height := 30.0 + rand.Float64()*50.0 // Random height between 30-80

	obstacle := &Obstacle{
		X:      rand.Float64() * (constants.WorldWidth - width),
		Y:      rand.Float64() * (constants.WorldHeight - height),
		Width:  width,
		Height: height,
		Active: true,
	}

	op.obstacles = append(op.obstacles, obstacle)
}

func (op *ObstaclePlugin) CheckCollision(x, y float64) bool {
	for _, obstacle := range op.obstacles {
		if obstacle.Active {
			if x >= obstacle.X && x <= obstacle.X+obstacle.Width &&
				y >= obstacle.Y && y <= obstacle.Y+obstacle.Height {
				return true
			}
		}
	}
	return false
}

func (op *ObstaclePlugin) GetObstacles() []*Obstacle {
	return op.obstacles
}

func (op *ObstaclePlugin) CheckCollisionRect(x, y, width, height float64) bool {
	for _, obstacle := range op.obstacles {
		if obstacle.Active {
			// Check for rectangle overlap
			if x < obstacle.X+obstacle.Width &&
				x+width > obstacle.X &&
				y < obstacle.Y+obstacle.Height &&
				y+height > obstacle.Y {
				return true
			}
		}
	}
	return false
}
