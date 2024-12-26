package stats

import (
	"fmt"
	"game/internal/core"
	"game/internal/plugins"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type StatsPlugin struct {
	kernel       *core.GameKernel
	playerPlugin plugins.PlayerPlugin
	gameFont     font.Face
}

func NewStatsPlugin(playerPlugin plugins.PlayerPlugin) *StatsPlugin {
	return &StatsPlugin{
		playerPlugin: playerPlugin,
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

	sp.gameFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (sp *StatsPlugin) Update() error {
	// No update logic needed for stats
	return nil
}

func (sp *StatsPlugin) Draw(screen *ebiten.Image) {
	playerPower := sp.playerPlugin.GetHealth()
	text.Draw(screen, fmt.Sprintf("Life: %.0f", playerPower), sp.gameFont, 10, 30, color.White)

	playerLevel := sp.playerPlugin.GetLevel()
	text.Draw(screen, fmt.Sprintf("Level: %.0f", playerLevel), sp.gameFont, 10, 60, color.White)

	playerNextLevelPercentage := sp.playerPlugin.NextLevelPercentage()
	percentageText := fmt.Sprintf("Next Level: %.0f%%", playerNextLevelPercentage*100)
	text.Draw(screen, percentageText, sp.gameFont, 10, 90, color.White)

	playerGetArmor := sp.playerPlugin.GetArmor()
	armorText := fmt.Sprintf("Armor: %.0f%%", playerGetArmor)
	text.Draw(screen, armorText, sp.gameFont, 10, 120, color.White)

	playerGetDamagePercent := sp.playerPlugin.GetDamagePercent()
	damagePercentage := fmt.Sprintf("Damage Percentage: %.0f%%", playerGetDamagePercent)
	text.Draw(screen, damagePercentage, sp.gameFont, 10, 150, color.White)

	playerGetSpeed := sp.playerPlugin.GetSpeed()
	speedText := fmt.Sprintf("Speed: %.0f%%", playerGetSpeed)
	text.Draw(screen, speedText, sp.gameFont, 10, 180, color.White)

	playerGetHealthRegenRate := sp.playerPlugin.GetHealthRegenRate()
	healthRegenRateText := fmt.Sprintf("Health Regen Rate: %.0f%%", playerGetHealthRegenRate)
	text.Draw(screen, healthRegenRateText, sp.gameFont, 10, 210, color.White)

	playerGetHealthRegenDelay := sp.playerPlugin.GetHealthRegenDelay()
	healthRegenDelayText := fmt.Sprintf("Health Regen Delay: %.0f%%", playerGetHealthRegenDelay)
	text.Draw(screen, healthRegenDelayText, sp.gameFont, 10, 240, color.White)

	playerGetAdditionalDamagePercent := sp.playerPlugin.GetAdditionalDamagePercent()
	additionalDamagePercentText := fmt.Sprintf("Additional Damage Percent: %.0f%%", playerGetAdditionalDamagePercent)
	text.Draw(screen, additionalDamagePercentText, sp.gameFont, 10, 270, color.White)

	playerGetCriticalChance := sp.playerPlugin.GetCriticalChance()
	criticalChanceText := fmt.Sprintf("Critical Chance: %.0f%%", playerGetCriticalChance)
	text.Draw(screen, criticalChanceText, sp.gameFont, 10, 300, color.White)
}
