// internal/game/game.go
package game

import (
	"game/internal/constants"
	"game/internal/core"
	"game/internal/game/components/menu"
	"game/internal/game/components/playingstate"
	"game/internal/game/states"
	"game/internal/helpers/fontface"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type Game struct {
	kernel            *core.GameKernel
	currentState      states.State
	componentsByState map[states.State]states.GameState
}

func NewGame(kernel *core.GameKernel) *Game {
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	fontface.FontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	menuState := menu.NewComponentMenuState(kernel)
	playingState := playingstate.NewComponentPlayingState(kernel)

	return &Game{
		kernel:       kernel,
		currentState: states.MenuState,
		componentsByState: map[states.State]states.GameState{
			states.MenuState:    menuState,
			states.PlayingState: playingState,
		},
	}
}

func (g *Game) Update() error {
	g.componentsByState[g.currentState].Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.componentsByState[g.currentState].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}

func (g *Game) SetState(state states.State) {
	g.currentState = state
}
