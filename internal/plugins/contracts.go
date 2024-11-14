package plugins

import (
	"game/internal/core"
	enemyentity "game/internal/plugins/enemy/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyPlugin interface {
	ID() string
	Init(kernel *core.GameKernel) error
	Update() error
	Draw(screen *ebiten.Image)
	Spawn(x, y float64)
	GetEnemies() []*enemyentity.Enemy
}
