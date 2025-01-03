package player

import (
	"game/internal/assets"
	"game/internal/config"
	"game/internal/constants"
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
	damagePercent  float64
	criticalChance float64

	maxHealth        float64
	healthRegenRate  float64
	healthRegenDelay float64
	healthRegenTimer float64

	walkingLeftAnimation  *assets.Animation
	walkingRightAnimation *assets.Animation
	idleAnimation         *assets.Animation
	currentAnimation      *assets.Animation

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
	additionalDamagePercentPerLevel float64

	dashSpeed    float64
	dashDuration float64
	dashCooldown float64
	dashTimer    float64
	isDashing    bool
	canDash      bool
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

func NewPlayerPlugin(plugins *core.PluginManager, c entities.Character) *PlayerPlugin {
	return &PlayerPlugin{
		playingPlugins: plugins,
		x:              400,
		y:              300,
		width:          32,
		height:         48,
		speed:          c.Speed,
		experience:     0,
		level:          1,

		health:           c.Health,
		maxHealth:        c.Health,
		healthRegenRate:  5.0,
		healthRegenDelay: 1.0,
		healthRegenTimer: 0,

		additionalDamagePercent: 10,
		criticalMultiplier:      2.0,
		armor:                   0,

		healthIncrementPerLevel:         10.0,
		speedIncrementPerLevel:          1.0,
		damagePercentIncrementPerLevel:  5.0,
		armorIncrementPerLevel:          1.0,
		criticalChanceIncrementPerLevel: 0.5,
		additionalDamagePercentPerLevel: 2.0,

		dashSpeed:    500, // Dash speed multiplier
		dashDuration: 0.2, // How long dash lasts
		dashCooldown: 1.0, // Time between dashes
		dashTimer:    0,
		isDashing:    false,
		canDash:      true,
	}
}

func (p *PlayerPlugin) ID() string {
	return "PlayerSystem"
}

func (p *PlayerPlugin) Init(kernel *core.GameKernel) error {
	p.kernel = kernel

	walkingLeftAnimation := assets.NewAnimation(0.1)
	err := walkingLeftAnimation.LoadFromJSON(
		"assets/images/player/rogue/run/left/asset.json",
		"assets/images/player/rogue/run/left/asset.png")

	if err != nil {
		log.Fatal("Failed to load player asset left:", err)
	}

	p.walkingLeftAnimation = walkingLeftAnimation

	walkingRightAnimation := assets.NewAnimation(0.1)
	err = walkingRightAnimation.LoadFromJSON(
		"assets/images/player/rogue/run/right/asset.json",
		"assets/images/player/rogue/run/right/asset.png")

	if err != nil {
		log.Fatal("Failed to load player asset right:", err)
	}

	p.walkingRightAnimation = walkingRightAnimation

	idleAnimation := assets.NewAnimation(0.1)
	err = idleAnimation.LoadFromJSON(
		"assets/images/player/rogue/idle/asset.json",
		"assets/images/player/rogue/idle/asset.png")

	if err != nil {
		log.Fatal("Failed to load player asset right:", err)
	}

	p.idleAnimation = idleAnimation

	return nil
}

func (p *PlayerPlugin) Update() error {
	// Get initial position
	newX, newY := p.x, p.y

	// Handle dash input and state
	if ebiten.IsKeyPressed(ebiten.KeyShift) && p.canDash {
		p.isDashing = true
		p.canDash = false
		p.dashTimer = 0
	}

	// Update dash state
	if !p.canDash {
		p.dashTimer += p.kernel.DeltaTime
		if p.dashTimer >= p.dashCooldown {
			p.canDash = true
			p.dashTimer = 0
		}
	}

	// Calculate movement
	currentSpeed := p.speed
	if p.isDashing {
		currentSpeed = p.dashSpeed
		p.dashTimer += p.kernel.DeltaTime

		if p.dashTimer >= p.dashDuration {
			p.isDashing = false
		}
	}

	// Apply movement with dash speed if active
	newX, newY = InputHandler(p, newX, newY, currentSpeed)

	// Update animation states
	p.currentAnimation = p.idleAnimation

	if p.x != newX || p.y != newY {
		if p.x > newX {
			p.currentAnimation = p.walkingRightAnimation
		} else {
			p.currentAnimation = p.walkingLeftAnimation
		}
	}

	// Update animation
	if p.currentAnimation != nil {
		p.currentAnimation.Update(p.kernel.DeltaTime)
	}

	// Update position
	p.x, p.y = newX, newY

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
		vector.DrawFilledRect(
			screen,
			float32(screenX-p.width/2),
			float32(screenY-p.height/2),
			float32(p.width),
			float32(p.height),
			color.RGBA{255, 255, 0, 255},
			true)

		vector.DrawFilledRect(screen,
			0,
			0,
			constants.ScreenWidth,
			constants.ScreenHeight,
			color.RGBA{200, 0, 0, 2},
			true)

	} else {
		if config.IsDebugEnv() {
			vector.DrawFilledRect(
				screen,
				float32(screenX-p.width/2),
				float32(screenY-p.height/2),
				float32(p.width),
				float32(p.height),
				color.RGBA{255, 255, 0, 255},
				true)
		}

	}

	drawInput := assets.DrawInput{
		Width:  p.width,
		Height: p.height,
		X:      screenX - p.width/2,
		Y:      screenY - p.height/2,
	}

	if p.currentAnimation != nil {
		p.currentAnimation.Draw(screen, drawInput)
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

func (p *PlayerPlugin) ApplyDamage(damage float64) {
	// Aplicar a armadura para reduzir o dano
	effectiveDamage := damage * (1 - p.armor/100)
	p.DecreaseHealth(effectiveDamage)
}

func (p *PlayerPlugin) CalculateDamage(baseDamage float64) (float64, bool) {
	isCriticalDamage := false
	damage := baseDamage * (1 + p.additionalDamagePercent/100)

	if rand.Float64() < p.criticalChance/100 {
		damage *= p.criticalMultiplier
		isCriticalDamage = true
	}

	return damage, isCriticalDamage
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
	p.health = p.maxHealth
	p.speed += p.speedIncrementPerLevel
	p.damagePercent += p.damagePercentIncrementPerLevel
	p.armor += p.armorIncrementPerLevel
	p.criticalChance += p.criticalChanceIncrementPerLevel
	p.additionalDamagePercent += p.additionalDamagePercentPerLevel
}

func (p *PlayerPlugin) GetDashTimer() float64 {
	return p.dashTimer
}

func (p *PlayerPlugin) GetNextLevelExperience() int {
	return levelUpExperience[p.level]
}
