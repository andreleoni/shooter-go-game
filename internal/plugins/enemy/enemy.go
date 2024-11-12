package enemy

import (
	"game/internal/core"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Enemy struct {
	X, Y   float64
	Speed  float64
	Active bool
}

type EnemyPlugin struct {
	kernel     *core.GameKernel
	enemies    []*Enemy
	spawnTimer float64
}

func NewEnemyPlugin() *EnemyPlugin {
	return &EnemyPlugin{
		enemies:    []*Enemy{},
		spawnTimer: 0,
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

	for _, enemy := range ep.enemies {
		if enemy.Active {
			enemy.Y += enemy.Speed * ep.kernel.DeltaTime
			if enemy.Y > 600 {
				enemy.Active = false
			}
		}
	}
	return nil
}

func (ep *EnemyPlugin) Draw(screen *ebiten.Image) {
	for _, enemy := range ep.enemies {
		if enemy.Active {
			ebitenutil.DrawRect(screen, enemy.X, enemy.Y, 20, 20, color.RGBA{0, 0, 255, 255})
		}
	}
}

func (ep *EnemyPlugin) Spawn(x, y float64) {
	ep.enemies = append(ep.enemies, &Enemy{X: x, Y: y, Speed: 100, Active: true})
}

func (ep *EnemyPlugin) GetEnemies() []*Enemy {
	return ep.enemies
}
