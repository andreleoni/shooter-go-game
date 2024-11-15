package core

import (
	"fmt"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Plugin interface {
	ID() string
	Init(kernel *GameKernel) error
	Update() error
	Draw(screen *ebiten.Image)
}

type PluginManager struct {
	plugins map[string]RegisteredPlugin
}

type RegisteredPlugin struct {
	plugin   Plugin
	priority int
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]RegisteredPlugin),
	}
}

func (pm *PluginManager) Register(plugin Plugin, priority int) {
	registeredPlugin := RegisteredPlugin{plugin: plugin}
	registeredPlugin.priority = priority

	pm.plugins[plugin.ID()] = registeredPlugin
}

func (pm *PluginManager) UpdateAll() error {
	for _, plugin := range retrieveSortedPlugins(pm.plugins) {
		fmt.Println(plugin.plugin.ID())

		if err := plugin.plugin.Update(); err != nil {
			return err
		}
	}

	return nil
}

func (pm *PluginManager) DrawAll(screen *ebiten.Image) {
	time.Sleep(100 * time.Millisecond)

	for _, plugin := range retrieveSortedPlugins(pm.plugins) {
		plugin.plugin.Draw(screen)
	}
}

func (pm *PluginManager) GetPlugin(id string) Plugin {
	return pm.plugins[id].plugin
}

func retrieveSortedPlugins(plugins map[string]RegisteredPlugin) []RegisteredPlugin {
	var arrayPlugins []RegisteredPlugin

	for _, plugin := range plugins {
		arrayPlugins = append(arrayPlugins, plugin)
	}

	sort.Slice(arrayPlugins, func(i, j int) bool {
		return arrayPlugins[i].priority < arrayPlugins[j].priority
	})
	fmt.Println(arrayPlugins)

	return arrayPlugins
}
