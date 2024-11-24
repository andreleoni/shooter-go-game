package camera

import (
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	X, Y float64
}

type CameraPlugin struct {
	kernel *core.GameKernel
	camera *Camera
	target plugins.PlayerPlugin
}

func NewCameraPlugin(target plugins.PlayerPlugin) *CameraPlugin {
	return &CameraPlugin{
		camera: &Camera{},
		target: target,
	}
}

func (cp *CameraPlugin) ID() string {
	return "CameraSystem"
}

func (cp *CameraPlugin) Init(kernel *core.GameKernel) error {
	cp.kernel = kernel
	return nil
}

func (cp *CameraPlugin) Update() error {
	playerX, playerY := cp.target.GetPosition()

	// Center camera on player
	cp.camera.X = playerX - constants.ScreenWidth/2
	cp.camera.Y = playerY - constants.ScreenHeight/2

	// Clamp camera to world bounds
	cp.camera.X = math.Max(0, math.Min(cp.camera.X, constants.WorldWidth-constants.ScreenWidth))
	cp.camera.Y = math.Max(0, math.Min(cp.camera.Y, constants.WorldHeight-constants.ScreenHeight))

	return nil
}

func (cp *CameraPlugin) Draw(*ebiten.Image) {
	// Draw camera
}

func (cp *CameraPlugin) GetPosition() (float64, float64) {
	return cp.camera.X, cp.camera.Y
}
