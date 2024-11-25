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

	kernel.EventBus.Subscribe("GameOver", func(data interface{}) {
		pluginManager.UnregisterAll()

		menuPlugin := menu.NewMenuPlugin(kernel)

		pluginManager.Register(menuPlugin, 0)

		menuPlugin.Init(kernel)
	})

	menuPlugin.Init(kernel)

	return &ComponentMenuState{kernel: kernel, pluginManager: pluginManager}
}

func (cms *ComponentMenuState) Draw(screen *ebiten.Image) {
	cms.pluginManager.DrawAll(screen)
}

func (cps *ComponentMenuState) PluginManager() *core.PluginManager {
	return cps.pluginManager
}
