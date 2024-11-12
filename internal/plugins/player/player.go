// plugins/player.go
package player

import (
	"game/internal/core"
	"game/internal/plugins/bullet"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type PlayerPlugin struct {
	kernel *core.GameKernel
	x, y   float64
	speed  float64
}

func NewPlayerPlugin() *PlayerPlugin {
	return &PlayerPlugin{
		x:     400,
		y:     300,
		speed: 200,
	}
}

func (p *PlayerPlugin) ID() string {
	return "PlayerSystem"
}

func (p *PlayerPlugin) Init(kernel *core.GameKernel) error {
	p.kernel = kernel
	return nil
}

func (p *PlayerPlugin) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.y -= p.speed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.y += p.speed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.x -= p.speed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.x += p.speed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		bulletPlugin := p.kernel.PluginManager.GetPlugin("BulletSystem").(*bullet.BulletPlugin)
		bulletPlugin.Shoot(p.x, p.y)
	}
	return nil
}

func (p *PlayerPlugin) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, p.x, p.y, 10, 10, color.RGBA{255, 255, 0, 255})
}

func (p *PlayerPlugin) GetPosition() (float64, float64) {
	return p.x, p.y
}
