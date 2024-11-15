// internal/game/game.go
package game

import (
	"game/internal/constants"
	"game/internal/core"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type Game struct {
	kernel         *core.GameKernel
	currentState   State
	characters     []Character
	selectedChar   int
	gameFont       font.Face
	canTransition  bool // Prevent multiple transitions
	selectionDelay float64
}

type Character struct {
	Name   string
	Speed  float64
	Health float64
}

func NewGame(kernel *core.GameKernel) *Game {
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	gameFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Game{
		kernel:       kernel,
		currentState: MenuState,
		characters: []Character{
			{Name: "Character 1", Speed: 1, Health: 100},
			{Name: "Character 2", Speed: 2, Health: 200},
		},
		selectionDelay: 0,
		gameFont:       gameFont,
	}
}

func (g *Game) Update() error {
	if !g.canTransition && !ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.canTransition = true
	}

	switch g.currentState {
	case MenuState:
		if g.canTransition && ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.currentState = CharacterSelectState
			g.canTransition = false
		}

	case CharacterSelectState:
		g.kernel.Update()

		g.selectionDelay += g.kernel.DeltaTime

		if g.selectionDelay > 0.05 {
			if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
				g.selectedChar = (g.selectedChar + 1) % len(g.characters)
				g.selectionDelay = 0
			}

			if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
				g.selectedChar--

				if g.selectedChar < 0 {
					g.selectedChar = len(g.characters) - 1
				}

				g.selectionDelay = 0
			}

			if g.canTransition && ebiten.IsKeyPressed(ebiten.KeyEnter) {
				g.currentState = PlayingState
				g.canTransition = false

				g.kernel.EventBus.Publish("StartGame", g.characters[g.selectedChar])

				g.selectionDelay = 0
			}
		}

	case GameOverState:
		if g.canTransition && ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.currentState = MenuState
			g.canTransition = false
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.currentState {
	case MenuState:
		text.Draw(screen, "Press ENTER to Start", g.gameFont, 300, 200, color.White)

	case CharacterSelectState:
		text.Draw(screen, "Select Character:", g.gameFont, 300, 150, color.White)

		for i, char := range g.characters {
			col := color.White

			if i == g.selectedChar {
				col = color.Gray16{200}
			}

			text.Draw(screen, char.Name, g.gameFont, 300, 200+(i*30), col)
		}

	case GameOverState:
		text.Draw(screen, "Game Over", g.gameFont, 350, 200, color.White)
		text.Draw(screen, "Press ENTER to Restart", g.gameFont, 300, 250, color.White)

	case PlayingState:
		g.kernel.PluginManager.DrawAll(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}
