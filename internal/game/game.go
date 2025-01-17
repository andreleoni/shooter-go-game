// internal/game/game.go
package game

import (
	"fmt"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/game/components/menu"
	"game/internal/game/components/playingstate"
	"game/internal/game/states"
	"game/internal/plugins/menu/fontface"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type Game struct {
	kernel            *core.GameKernel
	currentState      states.State
	componentsByState map[states.State]states.GameState

	updateCount int
	drawCount   int
	perSec      time.Time
}

func NewGame(kernel *core.GameKernel) *Game {
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	fontface.FontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	menuState := menu.NewComponentMenuState(kernel)
	playingState := playingstate.NewComponentPlayingState(kernel)

	game := &Game{
		kernel:       kernel,
		currentState: states.MenuState,
		componentsByState: map[states.State]states.GameState{
			states.MenuState:    menuState,
			states.PlayingState: playingState,
		},
	}

	ObserveStateChanges(game)

	return game
}

func (g *Game) Update() error {
	g.kernel.Update(g.componentsByState[g.currentState].PluginManager())

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.componentsByState[g.currentState].Draw(screen)

	now := time.Now()
	g.updateCount++

	if now.Sub(g.perSec) >= time.Second {
		fmt.Printf("TPS: %.2f, FPS: %.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
		fmt.Printf("Update() was called in this sec: %d times", g.updateCount)
		fmt.Printf("Draw() was called in this sec: %d times\n\n", g.drawCount)

		g.updateCount = 0
		g.drawCount = 0
		g.perSec = now
	}

	fpstext := fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS())

	text.Draw(
		screen,
		fpstext,
		basicfont.Face7x13,
		constants.ScreenWidth-80,
		20,
		color.White,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}

func (g *Game) SetState(state states.State) {
	g.currentState = state
}

func ObserveStateChanges(g *Game) {
	g.kernel.EventBus.Subscribe("StartGame", func(data interface{}) {
		g.SetState(states.PlayingState)
	})

	g.kernel.EventBus.Subscribe("GameOver", func(data interface{}) {
		g.SetState(states.MenuState)
	})
}
