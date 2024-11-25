package states

import (
	"game/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type State int

const (
	MenuState State = iota
	PlayingState
)

type GameState interface {
	PluginManager() *core.PluginManager
	Draw(screen *ebiten.Image)
}
