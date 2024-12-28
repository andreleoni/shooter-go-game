package player

import (
	"game/internal/assets"
	"game/internal/core"
	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/player/entities"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PlayerPlugin struct {
	kernel         *core.GameKernel
	playingPlugins *core.PluginManager

	x, y float64

	health         float64
	width          float64
	height         float64
	speed          float64
	armorPercent   float64
	damagePercent  float64
	criticalChance float64

	maxHealth        float64
	healthRegenRate  float64
	healthRegenDelay float64
	healthRegenTimer float64

	animation    *assets.Animation
	staticsprite *assets.StaticSprite

	facingRight bool

	DamageFlashTime float64

	experience int
	level      int

	additionalDamagePercent float64
	criticalMultiplier      float64

	armor float64

	healthIncrementPerLevel         float64
	speedIncrementPerLevel          float64
	damagePercentIncrementPerLevel  float64
	armorIncrementPerLevel          float64
	criticalChanceIncrementPerLevel float64
}

var levelUpExperience = map[int]int{
	// 1:  10,
	// 2:  12,
	// 3:  15,
	// 4:  20,
	// 5:  27,
	// 6:  40,
	// 7:  60,
	// 8:  90,
	// 9:  130,
	// 10: 200,
	1: 1,
	2: 2,
	3: 3,
	4: 4,
}

func NewPlayerPlugin(plugins *core.PluginManager, c entities.Character) *PlayerPlugin {
	return &PlayerPlugin{
		playingPlugins: plugins,
		x:              400,
		y:              300,
		width:          32,
		height:         32,
		speed:          c.Speed,
		experience:     0,
		level:          1,

		health:           c.Health,
		maxHealth:        c.Health,
		healthRegenRate:  5.0,
		healthRegenDelay: 1.0,
		healthRegenTimer: 0,

		additionalDamagePercent: 0,
		criticalMultiplier:      2.0,
		armor:                   0,

		healthIncrementPerLevel:         10.0,
		speedIncrementPerLevel:          1.0,
		damagePercentIncrementPerLevel:  5.0,
		armorIncrementPerLevel:          1.0,
		criticalChanceIncrementPerLevel: 0.5,
	}
}

func (p *PlayerPlugin) ID() string {
	return "PlayerSystem"
}

func (p *PlayerPlugin) Init(
	kernel *core.GameKernel,
) error {

	p.kernel = kernel

	p.animation = assets.NewAnimation(0.1)
	err := p.animation.LoadFromJSON("assets/images/player/gunner/run/player.json", "assets/images/player/gunner/run/player.png")
	if err != nil {
		log.Fatal("Failed to load player asset:", err)
	}

	// p.staticsprite = assets.NewStaticSprite()
	// err := p.staticsprite.Load("assets/images/player/player.png")
	// if err != nil {
	// 	log.Fatal("Failed to load player asset:", err)
	// // }

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

	if p.health < p.maxHealth {
		if p.healthRegenTimer >= p.healthRegenDelay {
			p.health += p.healthRegenRate * p.kernel.DeltaTime
			if p.health > p.maxHealth {
				p.health = p.maxHealth
			}

		} else {
			p.healthRegenTimer += p.kernel.DeltaTime
		}
	} else {
		p.healthRegenTimer = 0
	}

	return nil
}

func (p *PlayerPlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := p.playingPlugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	screenX := p.x - cameraX
	screenY := p.y - cameraY

	if p.DamageFlashTime > 0 {
		vector.DrawFilledRect(
			screen,
			float32(screenX-p.width/2),
			float32(screenY-p.height/2),
			float32(p.width),
			float32(p.height),
			color.RGBA{255, 255, 0, 255},
			true)

	} else {
		vector.DrawFilledRect(
			screen,
			float32(screenX-p.width/2),
			float32(screenY-p.height/2),
			float32(p.width),
			float32(p.height),
			color.RGBA{255, 255, 0, 255},
			true)
	}

	if p.animation != nil {
		p.animation.Update(p.kernel.DeltaTime)

		p.animation.Draw(screen, assets.DrawInput{
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

	if p.experience >= levelUpExperience[p.level] &&
		p.level < len(levelUpExperience)+1 {

		p.experience = 0
		p.level++

		p.increaseAttributes()

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

func (p *PlayerPlugin) AddAdditionalDamagePercent(percent float64) {
	p.additionalDamagePercent += percent
}

func (p *PlayerPlugin) AddArmor(amount float64) {
	p.armor += amount
}

func (p *PlayerPlugin) ApplyDamage(damage float64) {
	// Aplicar a armadura para reduzir o dano
	effectiveDamage := damage * (1 - p.armor/100)
	p.DecreaseHealth(effectiveDamage)
}

func (p *PlayerPlugin) CalculateDamage(baseDamage float64) float64 {
	// Aplicar a porcentagem de dano adicional
	damage := baseDamage * (1 + p.additionalDamagePercent/100)

	// Verificar se o ataque é um crítico
	if rand.Float64() < p.criticalChance/100 {
		damage *= p.criticalMultiplier
	}

	return damage
}

func (p *PlayerPlugin) GetAdditionalDamagePercent() float64 {
	return p.additionalDamagePercent
}

func (p *PlayerPlugin) GetArmor() float64 {
	return p.armor
}

func (p *PlayerPlugin) GetSpeed() float64 {
	return p.speed
}

func (p *PlayerPlugin) GetCriticalChance() float64 {
	return p.criticalChance
}

func (p *PlayerPlugin) GetDamagePercent() float64 {
	return p.damagePercent
}

func (p *PlayerPlugin) GetMaxHealth() float64 {
	return p.maxHealth
}

func (p *PlayerPlugin) GetHealthRegenRate() float64 {
	return p.healthRegenRate
}

func (p *PlayerPlugin) GetHealthRegenDelay() float64 {
	return p.healthRegenDelay
}

func (p *PlayerPlugin) GetHealthRegenTimer() float64 {
	return p.healthRegenTimer
}

func (p *PlayerPlugin) increaseAttributes() {
	p.maxHealth += p.healthIncrementPerLevel
	p.health = p.maxHealth // Restaurar a saúde ao máximo ao subir de nível
	p.speed += p.speedIncrementPerLevel
	p.damagePercent += p.damagePercentIncrementPerLevel
	p.armor += p.armorIncrementPerLevel
	p.criticalChance += p.criticalChanceIncrementPerLevel
}
