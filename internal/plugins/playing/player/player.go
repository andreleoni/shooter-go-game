package player

import (
	"game/internal/assets"
	"game/internal/core"
	"game/internal/plugins/camera"
	"game/internal/plugins/obstacle"
	"game/internal/plugins/weapon"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerPlugin struct {
	kernel         *core.GameKernel
	playingPlugins *core.PluginManager
	health         float64
	x, y           float64
	width          float64
	height         float64
	speed          float64
	animation      *assets.Animation
	staticsprite   *assets.StaticSprite
	shootTimer     float64
	shootCooldown  float64
	facingRight    bool

	level int
}

func NewPlayerPlugin() *PlayerPlugin {
	return &PlayerPlugin{
		x:             400,
		y:             300,
		speed:         200.0,
		shootCooldown: 1.0, // 1 second between shots
		shootTimer:    0,
		width:         32,
		height:        32,
		health:        100,
	}
}

func (p *PlayerPlugin) ID() string {
	return "PlayerSystem"
}

func (p *PlayerPlugin) Init(kernel *core.GameKernel) error {
	p.kernel = kernel

	p.staticsprite = assets.NewStaticSprite()
	err := p.staticsprite.Load("assets/images/player/player.png")
	if err != nil {
		log.Fatal("Failed to load player asset:", err)
	}

	return nil
}

func (p *PlayerPlugin) Update() error {
	newX, newY := p.x, p.y

	InputHandler(p, newX, newY)

	obstaclePlugin := p.playingPlugins.GetPlugin("ObstacleSystem").(*obstacle.ObstaclePlugin)
	if !obstaclePlugin.CheckCollisionRect(newX, newY, 20, 20) {
		p.x, p.y = newX, newY
	}

	// Auto-shooting
	p.shootTimer += p.kernel.DeltaTime
	if p.shootTimer >= p.shootCooldown {
		weapons := p.playingPlugins.GetPlugin("WeaponSystem").(*weapon.WeaponPlugin)
		weapons.Shoot(p.x, p.y)
		p.shootTimer = 0

	}

	return nil
}

func (p *PlayerPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := p.playingPlugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	screenX := p.x - cameraX
	screenY := p.y - cameraY

	p.staticsprite.DrawWithSize(
		screen,
		screenX-p.width/2,
		screenY-p.height/2,
		p.width,
		p.height,
		false)

	// p.animation.Draw(
	// 	screen,
	// 	screenX-p.width/2,
	// 	screenY-p.height/2,
	// 	p.width,
	// 	p.height,
	// 	!p.facingRight)
}

func (p *PlayerPlugin) GetPosition() (float64, float64) {
	return p.x, p.y
}

func (p *PlayerPlugin) DecreaseHealth(amount float64) {
	p.health -= amount

	if p.health < 0 {
		p.health = 0
		p.kernel.EventBus.Publish("GameOver", nil)
	}
}

func (p *PlayerPlugin) GetSize() (float64, float64) {
	return p.width, p.height
}

func (p *PlayerPlugin) GetHealth() float64 {
	return p.health
}
