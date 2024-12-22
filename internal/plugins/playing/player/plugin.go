package player

import (
	"game/internal/assets"
	"game/internal/core"
	"game/internal/plugins/playing/camera"
	"image/color"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type PlayerPlugin struct {
	kernel          *core.GameKernel
	playingPlugins  *core.PluginManager
	health          float64
	x, y            float64
	width           float64
	height          float64
	speed           float64
	animation       *assets.Animation
	staticsprite    *assets.StaticSprite
	facingRight     bool
	DamageFlashTime float64

	experience int
	level      int
}

var levelUpExperience = map[int]int{
	1:  10,
	2:  12,
	3:  15,
	4:  20,
	5:  27,
	6:  40,
	7:  60,
	8:  90,
	9:  130,
	10: 200,
}

func NewPlayerPlugin(plugins *core.PluginManager) *PlayerPlugin {
	return &PlayerPlugin{
		x:              400,
		y:              300,
		speed:          200.0,
		width:          32,
		height:         32,
		health:         100,
		playingPlugins: plugins,
		experience:     0,
		level:          1,
	}
}

func (p *PlayerPlugin) ID() string {
	return "PlayerSystem"
}

func (p *PlayerPlugin) Init(
	kernel *core.GameKernel,
) error {

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
	newX, newY = InputHandler(p, newX, newY)

	p.x, p.y = newX, newY

	// Atualizar o temporizador de flash de dano
	if p.DamageFlashTime > 0 {
		p.DamageFlashTime -= p.kernel.DeltaTime
	}

	return nil
}

func (p *PlayerPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := p.playingPlugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	screenX := p.x - cameraX
	screenY := p.y - cameraY

	if p.DamageFlashTime > 0 {
		ebitenutil.DrawRect(
			screen,
			screenX-p.width/2,
			screenY-p.height/2,
			p.width,
			p.height,
			color.RGBA{255, 255, 0, 255})

	} else {
		ebitenutil.DrawRect(
			screen,
			screenX-p.width/2,
			screenY-p.height/2,
			p.width,
			p.height,
			color.RGBA{255, 255, 0, 255})

		p.staticsprite.Draw(screen, assets.DrawInput{
			Width:  p.width,
			Height: p.height,
			X:      screenX - p.width/2,
			Y:      screenY - p.height/2,
		})
	}
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

func (p *PlayerPlugin) AddExperience(amount int) {

	p.experience += amount

	if p.experience >= levelUpExperience[p.level] {
		p.experience = 0
		p.level++

		p.kernel.EventBus.Publish("ChoosingAbility", nil)
	}
}

func (p *PlayerPlugin) GetLevel() float64 {
	return float64(p.level)
}

func (p *PlayerPlugin) GetExperience() float64 {
	return float64(p.experience)
}

func (p *PlayerPlugin) NextLevelPercentage() float64 {
	return float64(p.experience) / float64(levelUpExperience[p.level])
}
