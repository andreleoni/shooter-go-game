// internal/plugins/obstacle/obstacle.go
package obstacle

import (
	"game/internal/core"
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
	return nil
}

func (op *ObstaclePlugin) Update() error {
	op.spawnTimer += op.kernel.DeltaTime
	if op.spawnTimer >= 3.0 { // Spawn every 3 seconds
		op.SpawnRandomObstacle()
		op.spawnTimer = 0
	}
	return nil
}

func (op *ObstaclePlugin) Draw(screen *ebiten.Image) {
	for _, obstacle := range op.obstacles {
		if obstacle.Active {
			ebitenutil.DrawRect(screen,
				obstacle.X,
				obstacle.Y,
				obstacle.Width,
				obstacle.Height,
				color.RGBA{100, 100, 100, 255})
		}
	}
}

func (op *ObstaclePlugin) SpawnRandomObstacle() {
	width := 30.0 + rand.Float64()*50.0  // Random width between 30-80
	height := 30.0 + rand.Float64()*50.0 // Random height between 30-80

	obstacle := &Obstacle{
		X:      rand.Float64() * (800 - width),
		Y:      rand.Float64() * (600 - height),
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
