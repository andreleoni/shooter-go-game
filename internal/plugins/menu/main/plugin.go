package menu

import (
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"

	menu "game/internal/plugins/menu"
	"game/internal/plugins/menu/fontface"
	"image/color"
	"log"

	playerentities "game/internal/plugins/playing/player/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MenuPlugin struct {
	kernel *core.GameKernel

	currentState menu.State

	characters     []playerentities.Character
	selectedChar   int
	canTransition  bool
	selectionDelay float64

	initialMenuAnimation *assets.Animation
	backgroundAnimation  *assets.Animation
}

func NewMenuPlugin(kernel *core.GameKernel) *MenuPlugin {
	var err error

	initialMenuAnimation := assets.NewAnimation(0.1)
	err = initialMenuAnimation.LoadFromJSON(
		"assets/images/menu/initial/asset.json",
		"assets/images/menu/initial/asset.png")
	if err != nil {
		log.Fatal("Failed to load initial asset menu:", err)
	}

	backgroundAnimation := assets.NewAnimation(0.1)
	err = backgroundAnimation.LoadFromJSON(
		"assets/images/menu/background/asset.json",
		"assets/images/menu/background/asset.png")
	if err != nil {
		log.Fatal("Failed to load background asset menu:", err)
	}

	return &MenuPlugin{
		kernel:       kernel,
		currentState: menu.MenuState,
		characters: []playerentities.Character{
			{Name: "Rogue", Speed: 100, Health: 100, Ability: "Basic"},
		},
		selectionDelay:       0,
		initialMenuAnimation: initialMenuAnimation,
		backgroundAnimation:  backgroundAnimation,
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
	if !m.canTransition && !ebiten.IsKeyPressed(ebiten.KeyEnter) {
		m.canTransition = true
	}

	m.selectionDelay += m.kernel.DeltaTime

	switch m.currentState {
	case menu.MenuState:
		if m.canTransition && ebiten.IsKeyPressed(ebiten.KeyEnter) {
			m.currentState = menu.CharacterSelectState
			m.canTransition = false
		}

	case menu.CharacterSelectState:
		if m.selectionDelay > 0.1 {
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				m.selectedChar = (m.selectedChar + 1) % len(m.characters)
				m.selectionDelay = 0
			} else if ebiten.IsKeyPressed(ebiten.KeyW) {
				m.selectedChar--

				if m.selectedChar < 0 {
					m.selectedChar = len(m.characters) - 1
				}

				m.selectionDelay = 0
			}

			if m.canTransition && ebiten.IsKeyPressed(ebiten.KeyEnter) {
				m.currentState = menu.PlayingState
				m.canTransition = false

				m.kernel.EventBus.Publish("StartGame", m.characters[m.selectedChar])

				m.selectionDelay = 0
			}
		}

	case menu.GameOverState:
		if m.canTransition && ebiten.IsKeyPressed(ebiten.KeyEnter) {
			m.currentState = menu.MenuState
			m.canTransition = false
		}
	}

	return nil
}

func (m *MenuPlugin) Draw(screen *ebiten.Image) {
	switch m.currentState {
	case menu.MenuState:
		m.initialMenuAnimation.Draw(screen,
			assets.DrawInput{
				Width:  constants.ScreenWidth,
				Height: constants.ScreenHeight,
				X:      0,
				Y:      0,
			},
		)

		text.Draw(screen,
			"Press ENTER to Start",
			fontface.FontFace,
			(constants.ScreenWidth/2)-110,
			477,
			color.White)

	case menu.CharacterSelectState:
		m.backgroundAnimation.Draw(screen,
			assets.DrawInput{
				Width:  constants.ScreenWidth,
				Height: constants.ScreenHeight,
				X:      0,
				Y:      0,
			},
		)

		width := 300
		height := 300

		// colocar retangulo no meio da tela
		vector.DrawFilledRect(
			screen,
			float32((constants.ScreenWidth-width)/2),
			float32((constants.ScreenHeight-height)/2-50),
			float32(width),
			float32(height),
			color.RGBA{46, 41, 48, 255},
			true,
		)

		text.Draw(screen, "Select Character:", fontface.FontFace, 300, 150, color.White)

		for i, char := range m.characters {
			col := color.White

			if i == m.selectedChar {
				col = color.Gray16{200}
			}

			text.Draw(screen, char.Name, fontface.FontFace, 300, 200+(i*30), col)
		}

	case menu.GameOverState:
		text.Draw(screen, "Game Over", fontface.FontFace, 350, 200, color.White)
		text.Draw(screen, "Press ENTER to Restart", fontface.FontFace, 300, 250, color.White)

	}
}
