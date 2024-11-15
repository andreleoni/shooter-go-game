package player

import (
	"game/internal/core"
	"game/internal/plugins/bullet"
	"game/internal/plugins/camera"
	"game/internal/plugins/obstacle"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type PlayerPlugin struct {
	kernel        *core.GameKernel
	x, y          float64
	width         float64
	height        float64
	speed         float64
	shootTimer    float64
	shootCooldown float64
}

func NewPlayerPlugin() *PlayerPlugin {
	return &PlayerPlugin{
		x:             400,
		y:             300,
		speed:         200,
		shootCooldown: 1.0, // 1 second between shots
		shootTimer:    0,
		width:         20,
		height:        20,
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
	newX, newY := p.x, p.y

	oldX, oldY := p.x, p.y

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		newY -= p.speed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		newY += p.speed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		newX -= p.speed * p.kernel.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		newX += p.speed * p.kernel.DeltaTime
	}

	obstaclePlugin := p.kernel.PluginManager.GetPlugin("ObstacleSystem").(*obstacle.ObstaclePlugin)
	if !obstaclePlugin.CheckCollisionRect(newX, newY, 20, 20) {
		p.x, p.y = newX, newY
	} else {
		p.x, p.y = oldX, oldY
	}

	// Auto-shooting
	p.shootTimer += p.kernel.DeltaTime
	if p.shootTimer >= p.shootCooldown {
		bulletPlugin := p.kernel.PluginManager.GetPlugin("BulletSystem").(*bullet.BulletPlugin)
		bulletPlugin.Shoot(p.x, p.y)
		p.shootTimer = 0
	}

	return nil
}

func (p *PlayerPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := p.kernel.PluginManager.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	screenX := p.x - cameraX
	screenY := p.y - cameraY

	ebitenutil.DrawRect(screen, screenX, screenY, p.width, p.height, color.RGBA{255, 255, 0, 255})
}

func (p *PlayerPlugin) GetPosition() (float64, float64) {
	return p.x, p.y
}
