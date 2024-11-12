// core/plugin_manager.go
package core

import "github.com/hajimehoshi/ebiten/v2"

type Plugin interface {
	ID() string
	Init(kernel *GameKernel) error
	Update() error
	Draw(screen *ebiten.Image)
}

type PluginManager struct {
	plugins map[string]Plugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

func (pm *PluginManager) Register(plugin Plugin) {
	pm.plugins[plugin.ID()] = plugin
}

func (pm *PluginManager) UpdateAll() error {
	for _, plugin := range pm.plugins {
		if err := plugin.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (pm *PluginManager) DrawAll(screen *ebiten.Image) {
	for _, plugin := range pm.plugins {
		plugin.Draw(screen)
	}
}
