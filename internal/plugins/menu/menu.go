package menu

import (
	"game/internal/core"
	"game/internal/core/gamestate"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type MenuPlugin struct {
	kernel            *core.GameKernel
	currentState      gamestate.State
	characters        []Character
	selectedCharacter int
	gameFont          font.Face
}

type Character struct {
	Name   string
	Speed  float64
	Health float64
}

func NewMenuPlugin() *MenuPlugin {
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

	return &MenuPlugin{
		characters:   []Character{{Name: "Character 1", Speed: 1, Health: 100}, {Name: "Character 2", Speed: 2, Health: 200}},
		currentState: gamestate.MenuState,
		gameFont:     gameFont,
	}
}

func (m *MenuPlugin) ID() string {
	return "MenuPlugin"
}

func (m *MenuPlugin) Init(kernel *core.GameKernel) error {
	m.kernel = kernel
	return nil
}

func (m *MenuPlugin) Update() error {
	switch m.currentState {
	case gamestate.MenuState:
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			m.currentState = gamestate.CharacterSelectState
		}

	case gamestate.CharacterSelectState:
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			m.selectedCharacter = (m.selectedCharacter + 1) % len(m.characters)
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			m.selectedCharacter--
			if m.selectedCharacter < 0 {
				m.selectedCharacter = len(m.characters) - 1
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			m.currentState = gamestate.PlayingState
		}

	case gamestate.GameOverState:
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			m.currentState = gamestate.MenuState
		}
	}
	return nil
}

func (m *MenuPlugin) Draw(screen *ebiten.Image) {
	switch m.currentState {
	case gamestate.MenuState:
		text.Draw(screen, "Press Enter to Start", m.gameFont, 300, 200, color.White)

	case gamestate.CharacterSelectState:
		text.Draw(screen, "Select Character:", m.gameFont, 300, 150, color.White)
		for i, char := range m.characters {
			textColor := color.White
			if i == m.selectedCharacter {
				textColor = color.Gray16{255}
			}
			text.Draw(screen, char.Name, m.gameFont, 300, 200+i*30, textColor)
		}

	case gamestate.GameOverState:
		text.Draw(screen, "Game Over", m.gameFont, 350, 200, color.White)
		text.Draw(screen, "Press Enter to Restart", m.gameFont, 300, 250, color.White)
	}
}
