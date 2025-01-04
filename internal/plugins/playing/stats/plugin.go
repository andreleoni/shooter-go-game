package stats

import (
	"fmt"
	"game/internal/assets"
	"game/internal/core"
	"game/internal/plugins"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"

	abilityplugin "game/internal/plugins/playing/ability"
	"game/internal/plugins/playing/camera"
)

type StatsPlugin struct {
	kernel         *core.GameKernel
	playingPlugins *core.PluginManager

	playerPlugin plugins.PlayerPlugin
	gameFont     font.Face

	healthBarAnimation *assets.Animation

	showStats     bool
	showStatsTime float64
}

func NewStatsPlugin(plugins *core.PluginManager) *StatsPlugin {
	return &StatsPlugin{
		playingPlugins: plugins,
	}
}

func (sp *StatsPlugin) ID() string {
	return "StatsSystem"
}

func (sp *StatsPlugin) Init(kernel *core.GameKernel) error {
	sp.kernel = kernel

	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	healthBarAnimation := assets.NewAnimation(0.1)
	err = healthBarAnimation.LoadFromJSON(
		"assets/images/stats/bar/asset.json",
		"assets/images/stats/bar/asset.png")
	if err != nil {
		log.Fatal("Failed to load initial asset menu:", err)
	}

	sp.healthBarAnimation = healthBarAnimation

	sp.gameFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (sp *StatsPlugin) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyTab) {
		sp.showStatsTime = 0.1
	}

	if sp.showStatsTime > 0 {
		sp.showStatsTime -= sp.kernel.DeltaTime
		sp.showStats = true
	} else {
		sp.showStats = false
	}

	return nil
}

func (sp *StatsPlugin) Draw(screen *ebiten.Image) {
	playerPlugin := sp.playingPlugins.GetPlugin("PlayerSystem").(plugins.PlayerPlugin)

	if sp.showStats {
		playerPower := playerPlugin.GetHealth()
		text.Draw(screen, fmt.Sprintf("Life: %.0f", playerPower), sp.gameFont, 10, 30, color.White)

		playerLevel := playerPlugin.GetLevel()
		text.Draw(screen, fmt.Sprintf("Level: %.0f", playerLevel), sp.gameFont, 10, 60, color.White)

		playerNextLevelPercentage := playerPlugin.NextLevelPercentage()
		percentageText := fmt.Sprintf("Next Level: %.0f%%", playerNextLevelPercentage*100)
		text.Draw(screen, percentageText, sp.gameFont, 10, 90, color.White)

		playerGetArmor := playerPlugin.GetArmor()
		armorText := fmt.Sprintf("Armor: %.0f%%", playerGetArmor)
		text.Draw(screen, armorText, sp.gameFont, 10, 120, color.White)

		playerGetDamagePercent := playerPlugin.GetDamagePercent()
		damagePercentage := fmt.Sprintf("Damage Percentage: %.0f%%", playerGetDamagePercent)
		text.Draw(screen, damagePercentage, sp.gameFont, 10, 150, color.White)

		playerGetSpeed := playerPlugin.GetSpeed()
		speedText := fmt.Sprintf("Speed: %.0f%%", playerGetSpeed)
		text.Draw(screen, speedText, sp.gameFont, 10, 180, color.White)

		playerGetHealthRegenRate := playerPlugin.GetHealthRegenRate()
		healthRegenRateText := fmt.Sprintf("Health Regen Rate: %.0f%%", playerGetHealthRegenRate)
		text.Draw(screen, healthRegenRateText, sp.gameFont, 10, 210, color.White)

		playerGetHealthRegenDelay := playerPlugin.GetHealthRegenDelay()
		healthRegenDelayText := fmt.Sprintf("Health Regen Delay: %.0f%%", playerGetHealthRegenDelay)
		text.Draw(screen, healthRegenDelayText, sp.gameFont, 10, 240, color.White)

		playerGetAdditionalDamagePercent := playerPlugin.GetAdditionalDamagePercent()
		additionalDamagePercentText := fmt.Sprintf("Additional Damage Percent: %.0f%%", playerGetAdditionalDamagePercent)
		text.Draw(screen, additionalDamagePercentText, sp.gameFont, 10, 270, color.White)

		playerGetCriticalChance := playerPlugin.GetCriticalChance()
		criticalChanceText := fmt.Sprintf("Critical Chance: %.0f%%", playerGetCriticalChance)
		text.Draw(screen, criticalChanceText, sp.gameFont, 10, 300, color.White)

		abilities := sp.playingPlugins.GetPlugin("AbilitySystem")

		playerAbilities := abilities.(*abilityplugin.AbilityPlugin).GetAcquiredAbilities()

		for i, ability := range playerAbilities {
			ability := fmt.Sprintf("Ability: %s, Level: %d", ability.ID(), ability.CurrentLevel())
			text.Draw(screen, ability, sp.gameFont, 10, 330+(i*30), color.White)
		}

		// exibir o timer para o proximo dash do player em millisegundos
		playerDashTimer := playerPlugin.GetDashTimer()
		dashTimerText := fmt.Sprintf("Dash Cooldown: %.2fms", playerDashTimer)
		text.Draw(screen, dashTimerText, sp.gameFont, 10, 330+(len(playerAbilities)*30), color.White)
	}

	currentHealth := playerPlugin.GetHealth()
	maxHealth := playerPlugin.GetMaxHealth()
	healthPercentage := currentHealth / maxHealth

	currentXP := float64(playerPlugin.GetExperience())
	nextLevelXP := float64(playerPlugin.GetNextLevelExperience())
	xpPercentage := currentXP / nextLevelXP

	playerX, playerY := playerPlugin.GetPosition()

	_, playerHeight := playerPlugin.GetSize()

	cameraPlugin := sp.playingPlugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	screenX := playerX - cameraX
	screenY := playerY - cameraY

	barWidth := 50.0
	barHeight := 7.0
	expbarheight := 3.0

	healthBarPosition := screenY - (playerHeight / 2) - 15
	centerX := screenX - barWidth/2

	vector.DrawFilledRect(
		screen,
		float32(centerX),
		float32(healthBarPosition),
		float32(barWidth),
		float32(barHeight),
		color.RGBA{100, 0, 0, 255},
		true,
	)

	vector.DrawFilledRect(
		screen,
		float32(centerX),
		float32(healthBarPosition),
		float32(barWidth*healthPercentage),
		float32(barHeight),
		color.RGBA{255, 0, 0, 255},
		true,
	)

	vector.DrawFilledRect(
		screen,
		float32(centerX),
		float32(healthBarPosition+barHeight),
		float32(barWidth),
		float32(expbarheight),
		color.RGBA{0, 0, 100, 255},
		true,
	)

	vector.DrawFilledRect(
		screen,
		float32(centerX),
		float32(healthBarPosition+barHeight),
		float32(barWidth*xpPercentage),
		float32(expbarheight),
		color.RGBA{0, 100, 255, 255},
		true,
	)

	sp.healthBarAnimation.Draw(screen,
		assets.DrawInput{
			Width:  barWidth,
			Height: barHeight + expbarheight,
			X:      centerX,
			Y:      healthBarPosition,
		},
	)
}
