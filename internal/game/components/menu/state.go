package menu

import (
	"game/internal/core"
	menu "game/internal/plugins/menu/main"

	"github.com/hajimehoshi/ebiten/v2"
)

type ComponentMenuState struct {
	kernel        *core.GameKernel
	pluginManager *core.PluginManager
}

func NewComponentMenuState(kernel *core.GameKernel) *ComponentMenuState {
	pluginManager := core.NewPluginManager()

	menuPlugin := menu.NewMenuPlugin(kernel)
	pluginManager.Register(menuPlugin, 0)

	menuPlugin.Init(kernel)

	return &ComponentMenuState{kernel: kernel, pluginManager: pluginManager}
}

func (cms *ComponentMenuState) Update() {
	cms.pluginManager.UpdateAll()
}

func (cms *ComponentMenuState) Draw(screen *ebiten.Image) {
	cms.pluginManager.DrawAll(screen)
}
