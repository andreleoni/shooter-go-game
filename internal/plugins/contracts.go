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
	Spawn()
	GetEnemies() []*enemyentity.Enemy
}

type PlayerPlugin interface {
	ID() string
	Init(kernel *core.GameKernel) error
	Update() error
	Draw(screen *ebiten.Image)
	GetPosition() (float64, float64)
	DecreaseHealth(float64)
	GetSize() (float64, float64)
}
