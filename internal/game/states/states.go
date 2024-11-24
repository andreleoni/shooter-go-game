package states

import "github.com/hajimehoshi/ebiten/v2"

type State int

const (
	MenuState State = iota
	PlayingState
)

type GameState interface {
	Update()
	Draw(screen *ebiten.Image)
}
