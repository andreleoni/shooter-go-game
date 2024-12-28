package scenario

import (
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/playing/camera"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScenarioPlugin struct {
	kernel  *core.GameKernel
	plugins *core.PluginManager

	mapAnimation *assets.Animation
}

func New(plugins *core.PluginManager) *ScenarioPlugin {
	return &ScenarioPlugin{
		plugins: plugins,
	}
}

func (wp *ScenarioPlugin) ID() string {
	return "MapSystem"
}

func (mp *ScenarioPlugin) Init(kernel *core.GameKernel) error {
	mp.kernel = kernel

	mapAnimation := assets.NewAnimation(0.1)
	err := mapAnimation.LoadFromJSON(
		"assets/images/maps/grass/asset.json",
		"assets/images/maps/grass/asset.png")
	if err != nil {
		log.Fatal("Failed to map load animation:", err)
	}

	mp.mapAnimation = mapAnimation

	return nil
}

func (wp *ScenarioPlugin) Update() error {
	return nil
}

func (mp *ScenarioPlugin) Draw(screen *ebiten.Image) {
	if mp.mapAnimation != nil {
		cameraPlugin := mp.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
		cameraX, cameraY := cameraPlugin.GetPosition()

		mp.mapAnimation.Update(mp.kernel.DeltaTime)

		mp.mapAnimation.Draw(screen, assets.DrawInput{
			Width:  constants.WorldWidth,
			Height: constants.WorldHeight,
			X:      -cameraX,
			Y:      -cameraY,
		})
	}
}
