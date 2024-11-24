package abilities

import (
	"game/internal/core"
	"game/internal/plugins"

	"github.com/hajimehoshi/ebiten"
)

type AbilitiesPlugin struct {
	kernel *core.GameKernel
}

func NewAbilitiesPlugin(target plugins.PlayerPlugin) *AbilitiesPlugin {
	return &AbilitiesPlugin{}
}

func (cp *AbilitiesPlugin) ID() string {
	return "AbilitiesPlugin"
}

func (cp *AbilitiesPlugin) Init(kernel *core.GameKernel) error {
	cp.kernel = kernel
	return nil
}

func (cp *AbilitiesPlugin) Update() error {
	return nil
}

func (cp *AbilitiesPlugin) Draw(*ebiten.Image) {
	// Draw camera
}
