package experience

import (
	"game/internal/core"
	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/player"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ExperiencePlugin struct {
	kernel   *core.GameKernel
	crystals []*Crystal
	plugins  *core.PluginManager
}

type Crystal struct {
	X, Y   float64
	Active bool
}

func NewExperiencePlugin(plugins *core.PluginManager) *ExperiencePlugin {
	return &ExperiencePlugin{
		crystals: []*Crystal{},
		plugins:  plugins,
	}
}

func (ep *ExperiencePlugin) ID() string {
	return "ExperienceSystem"
}

func (ep *ExperiencePlugin) Init(kernel *core.GameKernel) error {
	ep.kernel = kernel
	rand.Seed(time.Now().UnixNano())
	return nil
}

func (ep *ExperiencePlugin) Update() error {
	playerPlugin := ep.plugins.GetPlugin("PlayerSystem").(*player.PlayerPlugin)

	playerX, playerY := playerPlugin.GetPosition()
	playerWidth, playerHeight := playerPlugin.GetSize()

	for _, crystal := range ep.crystals {
		if crystal.Active && ep.checkCollisionWithPlayer(crystal, playerX, playerY, playerWidth, playerHeight) {
			crystal.Active = false
			playerPlugin.AddExperience(10) // Adicionar experiência ao jogador
		}
	}

	return nil
}

func (ep *ExperiencePlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := ep.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for _, crystal := range ep.crystals {
		if crystal.Active {
			screenX := crystal.X - cameraX
			screenY := crystal.Y - cameraY

			ebitenutil.DrawRect(screen,
				screenX,
				screenY,
				10,
				10,
				color.RGBA{0, 255, 255, 255})
		}
	}
}

func (ep *ExperiencePlugin) DropCrystal(x, y float64) {
	ep.crystals = append(ep.crystals, &Crystal{
		X:      x,
		Y:      y,
		Active: true,
	})
}

func (ep *ExperiencePlugin) checkCollisionWithPlayer(crystal *Crystal, playerX, playerY, playerWidth, playerHeight float64) bool {
	return crystal.X < playerX+playerWidth && crystal.X+10 > playerX &&
		crystal.Y < playerY+playerHeight && crystal.Y+10 > playerY
}
